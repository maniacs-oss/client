'use strict'

import React, { Component, Text, TextInput, View, StyleSheet } from 'react-native'
import commonStyles from '../../styles/common'
import Button from '../../common-adapters/button'

export default class Account extends Component {

  constructor (props) {
    super(props)

    this.state = {
      email: this.props.email,
      oldPassphrase: null,
      newPassphrase: null,
      newPassphraseRepeat: null
    }
  }

  render () {
    const { email, emailVerified, onSave, emailError, passphraseError } = this.props
    return (
      <View style={styles.container}>
        <View style={styles.emailContainer}>
          <Text style={styles.verifiedTag}>
            {emailVerified ? 'Verified ✔' : 'Not Verified'}
          </Text>
          <TextInput style={commonStyles.textInput}
            onChangeText={(email) => this.setState({email})}
            defaultValue={email}/>
          {emailError && <Text style={[styles.errorInfo, {marginHorizontal: 10}]}>{emailError}</Text>}
        </View>

        <View style={styles.changePasswordContainer}>
          <Text style={{fontSize: 23, marginBottom: 20}}> Change Passphrase </Text>
          <TextInput style={commonStyles.textInput}
            returnKeyType='next'
            onSubmitEditing={(event) => this.refs['newPassphrase'].focus()}
            onChangeText={(oldPassphrase) => this.setState({oldPassphrase})}
            placeholder='Current passphrase'/>
          <TextInput style={commonStyles.textInput}
            returnKeyType='next'
            ref='newPassphrase'
            onSubmitEditing={(event) => this.refs['newPassphraseRepeat'].focus()}
            onChangeText={(newPassphrase) => this.setState({newPassphrase})}
            placeholder='New passphrase'/>
          <TextInput style={commonStyles.textInput}
            returnKeyType='next'
            ref='newPassphraseRepeat'
            onChangeText={(newPassphraseRepeat) => this.setState({newPassphraseRepeat})}
            placeholder='Confirm new passphrase'/>

          <View style={styles.saveContainer}>
            {passphraseError && <Text style={styles.errorInfo}>{passphraseError}</Text>}
            <View style={{flex: 1}}>
              <Button
                buttonStyle={commonStyles.button}
                style={styles.saveButton}
                title='Save'
                onPress={() => {
                  const { email, oldPassphrase, newPassphrase, newPassphraseRepeat } = this.state
                  onSave(email, oldPassphrase, newPassphrase, newPassphraseRepeat)
                }}/>
            </View>
          </View>
        </View>
      </View>
    )
  }

  static parseRoute () {
    return {
      componentAtTop: {
        title: 'Account',
        // Dummy data
        props: {
          email: 'kb-dawg@keybase.io',
          emailVerified: true,
          onSave: (email, oldPassphrase, newPassphrase, newPassphraseRepeat) => {
            console.log('saved! email:', email)
          }
        }
      }
    }
  }
}

Account.propTypes = {
  email: React.PropTypes.string.isRequired,
  emailVerified: React.PropTypes.bool.isRequired,
  onSave: React.PropTypes.func.isRequired,
  passphraseError: React.PropTypes.string,
  emailError: React.PropTypes.string
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    paddingHorizontal: 20
  },
  emailContainer: {
    flex: 0,
    marginBottom: 50,
    flexDirection: 'column'
  },
  changePasswordContainer: {
    flex: 1,
    flexDirection: 'column'
  },
  verifiedTag: {
    fontSize: 13,
    textAlign: 'right'
  },
  saveContainer: {
    flex: 0,
    marginTop: 14,
    marginHorizontal: 10,
    alignItems: 'center',
    flexDirection: 'row'
  },
  errorInfo: {
    color: 'red'
  },
  saveButton: {
    width: 80,
    alignSelf: 'flex-end'
  }
})
