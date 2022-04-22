import { Currency } from 'moonbeamswapdada'
import React, { useCallback, useEffect, useState } from 'react'
import ReactGA from 'react-ga'
import useLast from '../../hooks/useLast'
import { useSelectedListUrl } from '../../state/lists/hooks'
import { selectList, acceptListUpdate } from '../../state/lists/actions'
import Modal from '../Modal'
import { useDispatch } from 'react-redux'
import { AppDispatch } from '../../state'
import { CurrencySearch } from './CurrencySearch'
import ListIntroduction from './ListIntroduction'
import { ListSelect } from './ListSelect'
import { DEFAULT_TOKEN_LIST_URL } from '../../../src/constants/lists'
import Loader from 'react-loader-spinner'

interface CurrencySearchModalProps {
  isOpen: boolean
  onDismiss: () => void
  selectedCurrency?: Currency | null
  onCurrencySelect: (currency: Currency) => void
  otherSelectedCurrency?: Currency | null
  showCommonBases?: boolean
}

export default function CurrencySearchModal({
  isOpen,
  onDismiss,
  onCurrencySelect,
  selectedCurrency,
  otherSelectedCurrency,
  showCommonBases = false
}: CurrencySearchModalProps) {
  const [listView, setListView] = useState<boolean>(false)
  const [loading, setLoading] = useState<boolean>(false)
  const lastOpen = useLast(isOpen)

  useEffect(() => {
    if (isOpen && !lastOpen) {
      setListView(false)
    }
  }, [isOpen, lastOpen])

  const handleCurrencySelect = useCallback(
    (currency: Currency) => {
      onCurrencySelect(currency)
      onDismiss()
    },
    [onDismiss, onCurrencySelect]
  )

  const handleClickChangeList = useCallback(() => {
    ReactGA.event({
      category: 'Lists',
      action: 'Change Lists'
    })
    setListView(true)
  }, [])
  const handleClickBack = useCallback(() => {
    ReactGA.event({
      category: 'Lists',
      action: 'Back'
    })
    setListView(false)
  }, [])
  const listUrl = DEFAULT_TOKEN_LIST_URL
  const dispatch = useDispatch<AppDispatch>()
  const isSelected = listUrl === useSelectedListUrl()
  const handleSelectListIntroduction = useCallback(() => {
    if (isSelected) return
    ReactGA.event({
      category: 'Lists',
      action: 'Select List',
      label: listUrl
    })

    dispatch(selectList(listUrl))
    dispatch(acceptListUpdate(listUrl))
    setTimeout(function() {
      window.location.reload()
    }, 1000)
    setLoading(true)
  }, [dispatch, isSelected, listUrl])

  const selectedListUrl = useSelectedListUrl()
  const noListSelected = !selectedListUrl
  const loaderStyle = { display: 'flex', alignItems: 'center', marginLeft: 'auto', marginRight: 'auto' }

  return (
    <Modal isOpen={isOpen} onDismiss={onDismiss} maxHeight={90} minHeight={listView ? 40 : noListSelected ? 0 : loading ? 20 : 80}>
      {listView ? (
        <ListSelect onDismiss={onDismiss} onBack={handleClickBack} />
      ) : noListSelected ? (
        <ListIntroduction onSelectList={handleSelectListIntroduction} />
      ) : loading ? (
        <div style={loaderStyle}>
          <Loader
            type="Oval"
            color="#00BFFF"
            height={50}
            width={50}
            timeout={1000} // 3secs
          />
        </div>
      ) : (
        <CurrencySearch
          isOpen={isOpen}
          onDismiss={onDismiss}
          onCurrencySelect={handleCurrencySelect}
          onChangeList={handleClickChangeList}
          selectedCurrency={selectedCurrency}
          otherSelectedCurrency={otherSelectedCurrency}
          showCommonBases={showCommonBases}
        />
      )}
    </Modal>
  )
}
