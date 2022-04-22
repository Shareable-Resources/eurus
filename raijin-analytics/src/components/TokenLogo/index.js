import React, { useState, useEffect } from 'react'
import styled from 'styled-components'
import { isAddress } from '../../utils/index.js'
import EthereumLogo from '../../assets/eth.png'
import axios from 'axios'

const BAD_IMAGES = {}

const Inline = styled.div`
  display: flex;
  align-items: center;
  align-self: center;
`

const Image = styled.img`
  width: ${({ size }) => size};
  height: ${({ size }) => size};
  background-color: white;
  border-radius: 50%;
  box-shadow: 0px 6px 10px rgba(0, 0, 0, 0.075);
`

const StyledEthereumLogo = styled.div`
  display: flex;
  align-items: center;
  justify-content: center;

  > img {
    width: ${({ size }) => size};
    height: ${({ size }) => size};
  }
`

export default function TokenLogo({ address, symbol, header = false, size = '24px', ...rest }) {
  const [error, setError] = useState(false)
  // const [tokenInfo, setTokenInfo] = useState([])
  // const [tokenInfoArr, setTokenInfoArr] = useState([])
  // const [tokenLogo, setTokenLogo] = useState([])
  // const [tokenSymbol, setTokenSymbol] = useState([])

  // useEffect(() => {
  //   function fetchRaijinTokens() {
  //     axios
  //       .get('https://raw.githubusercontent.com/L-for-Louis/test/master/tokens.json')
  //       .then((response) => {
  //         let tokens = response.data.tokens
  //         setTokenInfo(tokens)
  //         // for (let i = 0; i < response.data.tokens.length; i++) {
  //         //   let tokens = response.data.tokens[i]
  //         //   tokenInfoArr.push(tokens)
  //         //   setTokenInfo(tokens)
  //         //   console.log(tokenInfoArr)
  //         // }
  //         // for (let i = 0; i < response.data.tokens.length; i++) {
  //         //   let logoURI = response.data.tokens[i].logoURI
  //         //   setTokenLogo(logoURI)
  //         //   return tokenInfo
  //         // }
  //         // for (let i = 0; i < response.data.tokens.length; i++) {
  //         //   let symbol = response.data.tokens[i].symbol
  //         //   setTokenSymbol(symbol)
  //         // }
  //       })
  //       .catch((error) => {
  //         console.log(`ERROR: ${error}`)
  //       })
  //   }
  //   fetchRaijinTokens()
  // }, [])

  // console.log("tokenSymbol", tokenSymbol, "tokenLogo", tokenLogo)
  // console.log("tokenInfo", tokenInfo)

  useEffect(() => {
    setError(false)
  }, [address])
  if (error || BAD_IMAGES[address]) {
    return (
      <Inline>
        <span {...rest} alt={''} style={{ fontSize: size }} role="img" aria-label="face">
          🤔
        </span>
      </Inline>
    )
  }

  // hard coded fixes for trust wallet api issues
  if (address?.toLowerCase() === '0x5e74c9036fb86bd7ecdcb084a0673efc32ea31cb') {
    address = '0x42456d7084eacf4083f1140d3229471bba2949a8'
  }

  if (address?.toLowerCase() === '0xc011a73ee8576fb46f5e1c5751ca3b9fe0af2a6f') {
    address = '0xc011a72400e58ecd99ee497cf89e3775d4bd732f'
  }

  if (address?.toLowerCase() === '0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2') {
    return (
      <StyledEthereumLogo size={size} {...rest}>
        <img
          src={EthereumLogo}
          style={{
            boxShadow: '0px 6px 10px rgba(0, 0, 0, 0.075)',
            borderRadius: '24px',
          }}
          alt=""
        />
      </StyledEthereumLogo>
    )
  }

  // if (symbol === tokenSymbol) {
  //       {
  //         tokenInfo.map((item, index) => (
  //             <Inline key={index}>
  //               <img
  //                 src={item.logoURI}
  //                 alt={item.name}
  //               />
  //             </Inline>
  //         ))
  //       }

  // return (
  //     <Inline>
  //       <Image src={tokenLogo} size={size} />
  //     </Inline>
  // )
  // }

  // for (let i = 0; i < tokenInfo.length; i++) {
  //   console.log("logoURI", tokenInfo[i].logoURI)
  //   if (symbol === tokenInfo[i].symbol || address === tokenInfo[i].address) {
  //   // console.log("symbol", symbol, "address", address, "logoURI", tokenInfo.logoURI)
  //   return (
  //     <Inline>
  //       <Image src={tokenInfo[i].logoURI} size={size} />
  //     </Inline>
  //   )
  // }
  // }

  // console.log(tokenInfo.symbol, tokenInfo.address, tokenInfo.logoURI)

  // if (symbol === tokenInfo.symbol || address === tokenInfo.address) {
  //   // console.log("symbol", symbol, "address", address, "logoURI", tokenInfo.logoURI)
  //   return (
  //     <Inline>
  //       <Image src={tokenInfo.logoURI} size={size} />
  //     </Inline>
  //   )
  // }

  // if (symbol === tokenInfo.symbol || address === tokenInfo.address) {
  //   // console.log("symbol", symbol, "address", address, "logoURI", tokenInfo.logoURI)
  //   return (
  //     <Inline>
  //       <Image src={tokenInfo.logoURI} size={size} />
  //     </Inline>
  //   )
  // }
  // if (symbol === tokenInfo.symbol) {
  //   return (
  //   {tokenInfoArr.map((item, index) => (
  //     <Inline key={index}>
  //       <img
  //         src={item.logoURI}
  //         alt={item.name}
  //       />
  //     </Inline>
  //   ))}
  //   )
  // }

  // if (tokenSymbol.length > 0 ) {
  //   if (tokenSymbol.find(symbol)) {
  //     console.log(symbol)
  //   }
  // }

  if (symbol === 'EUN' || address.toLowerCase() === '0x17deba6e45745d6b72f684a150566001283e4424') {
    return (
      <Inline>
        <Image src={'https://raw.githubusercontent.com/L-for-Louis/token-list/main/images/WEUN.svg'} size={size} />
      </Inline>
    )
  } else if (symbol === 'WEUN' || address.toLowerCase() === '0xe39a0ea2f9a7c296f581c0c2e7cf8d52f96b4de0') {
    return (
      <Inline>
        <Image src={'https://raw.githubusercontent.com/L-for-Louis/token-list/main/images/WEUN.svg'} size={size} />
      </Inline>
    )
  } else if (symbol === 'ETH' || address.toLowerCase() === '0x166f79582d5127f878426a1876c8ae98fd25fd61') {
    return (
      <Inline>
        <Image src={'https://raw.githubusercontent.com/L-for-Louis/token-list/main/images/ETH.svg'} size={size} />
      </Inline>
    )
  } else if (symbol === 'DAI' || address.toLowerCase() === '0x84479a80707d746953ef24dc24cf8820116c3e8c') {
    return (
      <Inline>
        <Image
          src={
            'https://raw.githubusercontent.com/uniswap/assets/master/blockchains/ethereum/assets/0x6B175474E89094C44Da98b954EedeAC495271d0F/logo.png'
          }
          size={size}
        />
      </Inline>
    )
  } else if (symbol === 'USDT' || address.toLowerCase() === '0xa54dee79c3bb34251debf86c1ba7d21898ffb7ac') {
    return (
      <Inline>
        <Image
          alt={''}
          // src={path}
          src={'https://raw.githubusercontent.com/L-for-Louis/token-list/main/images/USDT.svg'}
          size={size}
        />
      </Inline>
    )
  } else if (symbol === 'USDC' || address.toLowerCase() === '0xc56743422c75098c2e2ebc2b78351e3702d5eadd') {
    return (
      <Inline>
        <Image src={'https://raw.githubusercontent.com/L-for-Louis/token-list/main/images/USDC.png'} size={size} />
      </Inline>
    )
  } else if (symbol === 'CDAI' || address.toLowerCase() === '0xe9c5aeb38a7e8762116cbe4a281a20ac12f89977') {
    return (
      <Inline>
        <Image
          src={'https://raw.githubusercontent.com/compound-finance/token-list/master/assets/ctoken_dai.svg'}
          size={size}
        />
      </Inline>
    )
  } else if (symbol === 'COMP' || address.toLowerCase() === '0x80028767e4925e4e20183dc76ab070a62985facb') {
    return (
      <Inline>
        <Image
          src={'https://raw.githubusercontent.com/compound-finance/token-list/master/assets/asset_COMP.svg'}
          size={size}
        />
      </Inline>
    )
  } else if (symbol === 'MKR' || address.toLowerCase() === '0x0a314f415f8eb4e58ccf962f1b559eb5562c058c') {
    return (
      <Inline>
        <Image src={'https://www.gemini.com/images/currencies/icons/default/mkr.svg'} size={size} />
      </Inline>
    )
  } else if (symbol === 'BAT' || address.toLowerCase() === '0x087d78cf6138c68fb660deb5e0fb498f31fc5cf8') {
    return (
      <Inline>
        <Image
          src={
            'https://www.google.com/url?sa=i&url=https%3A%2F%2Fwww.pikpng.com%2Fpngvi%2FibTmhTh_basic-attention-token-logo-bat-token-clipart%2F&psig=AOvVaw0tTxD3eqkS1z-3K9WZDuvT&ust=1639539535739000&source=images&cd=vfe&ved=0CAsQjRxqFwoTCOCn-e6u4vQCFQAAAAAdAAAAABAD'
          }
          size={size}
        />
      </Inline>
    )
  } else if (symbol === 'GAO' || address.toLowerCase() === '0x0a0850c65a8c47dc4133fe73795fa7cd92d46b71') {
    return (
      <Inline>
        <Image src={'https://static.wikia.nocookie.net/evchk/images/9/9c/HKGolden_Plastic_Icon.svg'} size={size} />
      </Inline>
    )
  } else {
    const path = `https://raw.githubusercontent.com/trustwallet/assets/master/blockchains/ethereum/assets/${isAddress(
      address
    )}/logo.png`

    return (
      <Inline>
        <Image
          {...rest}
          alt={''}
          src={path}
          size={size}
          onError={(event) => {
            BAD_IMAGES[address] = true
            setError(true)
            event.preventDefault()
          }}
        />
      </Inline>
    )
  }
}
