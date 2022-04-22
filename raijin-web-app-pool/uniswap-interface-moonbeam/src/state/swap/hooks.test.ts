import { parse } from 'qs'
import { Field } from './actions'
import { queryParametersToSwapState } from './hooks'

describe('hooks', () => {
  describe('#queryParametersToSwapState', () => {
    test('ETH to DAI', () => {
      expect(
        queryParametersToSwapState(
          parse(
            '?inputCurrency=ETH&outputCurrency=0xE39A0EA2F9a7C296F581c0C2e7Cf8d52F96B4De0&exactAmount=20.5&exactField=outPUT',
            { parseArrays: false, ignoreQueryPrefix: true }
          )
        )
      ).toEqual({
        [Field.OUTPUT]: { currencyId: '0xE39A0EA2F9a7C296F581c0C2e7Cf8d52F96B4De0' },
        [Field.INPUT]: { currencyId: 'EUN' },
        typedValue: '20.5',
        independentField: Field.OUTPUT,
        recipient: null
      })
    })

    test('does not duplicate eth for invalid output token', () => {
      expect(
        queryParametersToSwapState(parse('?outputCurrency=invalid', { parseArrays: false, ignoreQueryPrefix: true }))
      ).toEqual({
        [Field.INPUT]: { currencyId: '' },
        [Field.OUTPUT]: { currencyId: 'ETH' },
        typedValue: '',
        independentField: Field.INPUT,
        recipient: null
      })
    })

    test('output ETH only', () => {
      expect(
        queryParametersToSwapState(
          parse('?outputCurrency=eth&exactAmount=20.5', { parseArrays: false, ignoreQueryPrefix: true })
        )
      ).toEqual({
        [Field.OUTPUT]: { currencyId: 'ETH' },
        [Field.INPUT]: { currencyId: '' },
        typedValue: '20.5',
        independentField: Field.INPUT,
        recipient: null
      })
    })

    test('invalid recipient', () => {
      expect(
        queryParametersToSwapState(
          parse('?outputCurrency=eth&exactAmount=20.5&recipient=abc', { parseArrays: false, ignoreQueryPrefix: true })
        )
      ).toEqual({
        [Field.OUTPUT]: { currencyId: 'ETH' },
        [Field.INPUT]: { currencyId: '' },
        typedValue: '20.5',
        independentField: Field.INPUT,
        recipient: null
      })
    })

    test('valid recipient', () => {
      expect(
        queryParametersToSwapState(
          parse('?outputCurrency=eth&exactAmount=20.5&recipient=0x0fF2D1eFd7A57B7562b2bf27F3f37899dB27F4a5', {
            parseArrays: false,
            ignoreQueryPrefix: true
          })
        )
      ).toEqual({
        [Field.OUTPUT]: { currencyId: 'ETH' },
        [Field.INPUT]: { currencyId: '' },
        typedValue: '20.5',
        independentField: Field.INPUT,
        recipient: '0x0fF2D1eFd7A57B7562b2bf27F3f37899dB27F4a5'
      })
    })
    test('accepts any recipient', () => {
      expect(
        queryParametersToSwapState(
          parse('?outputCurrency=eth&exactAmount=20.5&recipient=bob.argent.xyz', {
            parseArrays: false,
            ignoreQueryPrefix: true
          })
        )
      ).toEqual({
        [Field.OUTPUT]: { currencyId: 'ETH' },
        [Field.INPUT]: { currencyId: '' },
        typedValue: '20.5',
        independentField: Field.INPUT,
        recipient: 'bob.argent.xyz'
      })
    })
  })
})
