// @flow

import React from 'react'
import {Box} from '../common-adapters'
import About from './about'
import SettingsNav from './nav/index.native'
import * as settingsConstants from '../constants/settings'

import type {DumbComponentMap} from '../constants/types/more'

const aboutMap: DumbComponentMap<About> = {
  component: About,
  mocks: {
    'Normal': {version: '1.0.18-20161107120015+aee424b.'},
  },
}

const settingsNavMap: DumbComponentMap<SettingsNav> = {
  component: SettingsNav,
  mocks: {
    'Normal': {
      selectedTab: settingsConstants.landingTab,
      onTabChange: (tab) => console.log('clicked', tab),
      badgeNumbers: {},
    },
  },
}

export default {
  About: aboutMap,
  SettingsNav: settingsNavMap,
}
