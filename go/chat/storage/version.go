package storage

import (
	"context"
	"fmt"

	"github.com/keybase/client/go/chat/utils"
	"github.com/keybase/client/go/libkb"
	"github.com/keybase/client/go/protocol/chat1"
)

type ServerVersions struct {
	libkb.Contextified
	utils.DebugLabeler

	cached *chat1.ServerCacheVers
}

func NewServerVersions(g *libkb.GlobalContext) *ServerVersions {
	return &ServerVersions{
		Contextified: libkb.NewContextified(g),
		DebugLabeler: utils.NewDebugLabeler(g, "ServerVersions", false),
	}
}

func (s *ServerVersions) makeKey() libkb.DbKey {
	return libkb.DbKey{
		Typ: libkb.DBChatBlocks,
		Key: fmt.Sprintf("versions"),
	}
}

func (s *ServerVersions) fetchLocked(ctx context.Context) (chat1.ServerCacheVers, Error) {
	// Check in memory first
	if s.cached != nil {
		return *s.cached, nil
	}

	// Read from LevelDb
	raw, found, err := s.G().LocalChatDb.GetRaw(s.makeKey())
	if err != nil {
		return chat1.ServerCacheVers{},
			NewInternalError(ctx, s.DebugLabeler, "GetRaw error: %s", err.Error())
	}
	if !found {
		s.Debug(ctx, "no server version found, using defaults")
		return chat1.ServerCacheVers{}, nil
	}
	var srvVers chat1.ServerCacheVers
	if err = decode(raw, &srvVers); err != nil {
		return chat1.ServerCacheVers{},
			NewInternalError(ctx, s.DebugLabeler, "decode error: %s", err.Error())
	}

	// Store in memory
	s.cached = &srvVers
	return *s.cached, nil
}

func (s *ServerVersions) Fetch(ctx context.Context) (chat1.ServerCacheVers, Error) {
	locks.Version.Lock()
	defer locks.Version.Unlock()

	return s.fetchLocked(ctx)
}

func (s *ServerVersions) matchLocked(ctx context.Context, vers int,
	versFunc func(chat1.ServerCacheVers) int) Error {
	srvVers, err := s.fetchLocked(ctx)
	if err != nil {
		return err
	}
	retVers := versFunc(srvVers)
	if retVers != vers {
		return NewVersionMismatchError(chat1.InboxVers(vers), chat1.InboxVers(retVers))
	}
	return nil
}

func (s *ServerVersions) MatchInbox(ctx context.Context, vers int) Error {
	locks.Version.Lock()
	defer locks.Version.Unlock()

	return s.matchLocked(ctx, vers, func(srvVers chat1.ServerCacheVers) int { return srvVers.InboxVers })
}

func (s *ServerVersions) MatchBodies(ctx context.Context, vers int) Error {
	locks.Version.Lock()
	defer locks.Version.Unlock()

	return s.matchLocked(ctx, vers, func(srvVers chat1.ServerCacheVers) int { return srvVers.BodiesVers })
}

func (s *ServerVersions) Sync(ctx context.Context, vers chat1.ServerCacheVers) Error {
	locks.Version.Lock()
	defer locks.Version.Unlock()

	// Write in memory
	s.cached = &vers

	// Write out to LevelDB
	dat, err := encode(vers)
	if err != nil {
		return NewInternalError(ctx, s.DebugLabeler, "encode error: %s", err.Error())
	}
	if err = s.G().LocalChatDb.PutRaw(s.makeKey(), dat); err != nil {
		return NewInternalError(ctx, s.DebugLabeler, "PutRaw error: %s", err.Error())
	}

	return nil
}
