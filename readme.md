# hello go bitcoin

## links
https://github.com/btcsuite/btcd

https://hou27.tistory.com/entry/Go%EB%A1%9C-%EB%A7%8C%EB%93%9C%EB%8A%94-%EB%B8%94%EB%A1%9D%EC%B2%B4%EC%9D%B8-part-5-Wallet


https://live.blockcypher.com/

https://github.com/btcsuite/btcd/blob/master/btcutil/address_test.go

https://github.com/miguelmota/bitcoin-development-with-go/blob/master/en/transfer-coin/README.md

### explorer

https://live.blockcypher.com/btc-testnet/

https://blockstream.info/testnet/

## 시작

```
go mod init github.com/ricepotato/hello-go-bitcoin
```

## 의존성 설치

x/crypto module 0.16 버전을 사용해야 함

```
go get -u github.com/btcsuite/btcd

go get golang.org/x/crypto@v0.16.0
```

## 읽을거리


### [비트코인 구조] 개인키(Private key), 공개키(Public key), 주소(Address) 생성
https://kwjdnjs.tistory.com/55


### [비트코인 구조] P2PKH
https://kwjdnjs.tistory.com/61

### [비트코인 구조] P2SH
https://kwjdnjs.tistory.com/62

### [비트코인 구조] 세그윗(Segwit), Bech32 주소
https://kwjdnjs.tistory.com/63

### 세그윗

https://medium.com/programming-bitcoin/chapter-13-%EC%84%B8%EA%B7%B8%EC%9C%97-865a0c3f6414

### 스크립트

https://medium.com/programming-bitcoin/chapter-6-%EC%8A%A4%ED%81%AC%EB%A6%BD%ED%8A%B8-2474f708091b

### 트랜잭션

https://medium.com/programming-bitcoin/chapter-5-%ED%8A%B8%EB%9E%9C%EC%9E%AD%EC%85%98-e5ced4ad04af

### 탭루트 슈노르 서명

https://kwjdnjs.tistory.com/86

https://kwjdnjs.tistory.com/87

https://academy.gopax.co.kr/taebruteuran-mueosimyeo-biteukoine-eoddeohge-doumi-doelggayo/

### 비트코인 주소 종류

https://blog.naver.com/wolsajubu/222937013012

네이티브 세그윗(bc1 ~ )과 세그윗 (3~)
https://www.ledger.com/ko/academy/%EC%84%B8%EA%B7%B8%EC%9C%97%EA%B3%BC-%EB%84%A4%EC%9D%B4%ED%8B%B0%EB%B8%8C-%EC%84%B8%EA%B7%B8%EC%9C%97bech32%EC%9D%80-%EC%96%B4%EB%96%A4-%EC%B0%A8%EC%9D%B4%EC%A0%90%EC%9D%B4-%EC%9E%88%EB%82%98

## token pocket privkey

token pocket 주소지원
https://medium.com/token-pocket-kr/%ED%86%A0%ED%81%B0%ED%8F%AC%EC%BC%93-4%EA%B0%80%EC%A7%80-%EC%9C%A0%ED%98%95%EC%9D%98-%EB%B9%84%ED%8A%B8%EC%BD%94%EC%9D%B8-%EC%A3%BC%EC%86%8C-%ED%98%95%EC%8B%9D-%EC%A7%80%EC%9B%90-3ae2edd362d2



## examples

### wallet generate

https://github.com/miguelmota/bitcoin-development-with-go/blob/master/en/wallet-generate/README.md

### transfer coin

https://github.com/miguelmota/bitcoin-development-with-go/blob/master/en/transfer-coin/README.md


## RPC 

https://www.quicknode.com/docs/bitcoin

## 오류들

### {"code":-26,"message":"dust"}

sendrawtransaction 호출 시 에러가 발생할 수 있다. dust 는 너무 적은 UTXO 를 생성하려 할 때 발생한다. 너무 적은 UTXO 는 그 가치보다 더 높은 수수료가 발생하기 때문에 애초에 생성되지 않는다.

https://bitcoin.stackexchange.com/questions/10986/what-is-meant-by-bitcoin-dust

```
Bitcoin Core considers a transaction output to be dust, when its value is lower than the cost of creating and spending it at the dustRelayFee rate. The default value for dustRelayFee is 3,000 sat/kvB¹, which results in the same dust values as the prior dust definition used before Bitcoin Core 0.15.0. The previous dust definition tied the dust limit to the minRelayTxFee rate and the spending cost of an output exceeding 1/3 of its value.
```


### not enough UTXOs to meet target amount

UTXO 가 부족할 때 발생하는 에러이다. bitcoin 잔액 부족.


### {"result":null,"error":{"code":-26,"message":"mandatory-script-verify-flag-failed (Witness requires empty scriptSig)"},"id":"1"}

legacy 방식으로 transaction 을 구성하여 native segwit 주소의 UTXO 를 소비하려 할때 sendrawtransaction 에서 발생하는 에러. 

native segwit 방식과 legacy 방식은 transaction 의 구성이 다르다. 아래 링크 참고.

https://medium.com/programming-bitcoin/chapter-13-%EC%84%B8%EA%B7%B8%EC%9C%97-865a0c3f6414

https://bitcoin.stackexchange.com/questions/77440/segwit-transaction-in-golang