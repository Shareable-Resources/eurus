import React from 'react'
import { Text } from 'rebass'
// import Logo from '../../assets/images/mainlogo.png'
import Logo from '../../assets/images/raijinLogo.png'
import { ButtonPrimary } from '../Button'
import Column, { AutoColumn } from '../Column'
import { PaddedColumn } from './styleds'


export default function ListIntroduction({ onSelectList }: { onSelectList: () => void }) {
  // const { t } = useTranslation()

  return (
    <Column style={{ width: '100%', flex: '1 1' }}>
    <PaddedColumn>
      <AutoColumn gap="14px">
        <img
          style={{ width: '120px', margin: '0 auto' }}
          src={Logo} 
          alt="token-list-preview"
        />
        <Text style={{ marginBottom: '8px', textAlign: 'center' }}>
          Eurus Swap supports a wide range token lists.
        </Text>
        <ButtonPrimary onClick={onSelectList} id="list-introduction-choose-a-list">
          Import Eurus Default
        </ButtonPrimary>
      </AutoColumn>
    </PaddedColumn>
  </Column>
  )
}
