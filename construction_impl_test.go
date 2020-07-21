/*******************************************************************************
*   (c) 2020 Zondax GmbH
*
*  Licensed under the Apache License, Version 2.0 (the "License");
*  you may not use this file except in compliance with the License.
*  You may obtain a copy of the License at
*
*      http://www.apache.org/licenses/LICENSE-2.0
*
*  Unless required by applicable law or agreed to in writing, software
*  distributed under the License is distributed on an "AS IS" BASIS,
*  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
*  See the License for the specific language governing permissions and
*  limitations under the License.
********************************************************************************/
package rosetta_filecoin_lib

import (
  "testing"
  "encoding/hex"
  "encoding/base64"
  "encoding/json"
  "reflect"
  "github.com/filecoin-project/lotus/chain/types"
  "github.com/filecoin-project/specs-actors/actors/abi"
)

func TestDeriveFromPublicKey(t *testing.T) {
  pk, err := hex.DecodeString("04fc016f3d88dc7070cdd95b5754d32fd5290f850b7c2208fca0f715d35861de1841d9a342a487692a63810a6c906b443a18aa804d9d508d69facc5b06789a01b4")
  if err != nil {
    t.Errorf("Invalid test case")
  }

  r := &RosettaConstructionFilecoin{false}

  address, err := r.DeriveFromPublicKey(pk)
  if err != nil {
    t.Errorf("FIX ME")
  }

  if address != "t1rovwtiuo5ncslpmpjftzu5akswbgsgighjazxoi" {
    t.Fail()
  }

}

func TestSign(t *testing.T) {
  unsignedTx := `{
    "To": "t17uoq6tp427uzv7fztkbsnn64iwotfrristwpryy",
    "From": "t1d2xrzcslx7xlbbylc5c3d5lvandqw4iwl6epxba",
    "Nonce": 1,
    "Value": "100000",
    "GasPrice": "2500",
    "GasLimit": 25000,
    "Method": 0,
    "Params": ""
  }`
  sk, err := hex.DecodeString("f15716d3b003b304b8055d9cc62e6b9c869d56cc930c3858d4d7c31f5f53f14a")
  if err != nil {
    t.Errorf("FIX ME")
  }
  r := &RosettaConstructionFilecoin{false}

  rawIn := json.RawMessage(unsignedTx)

  bytes, err := rawIn.MarshalJSON()
  if err != nil {
    t.Errorf("FIX ME")
  }

  var msg types.Message
  err = json.Unmarshal(bytes, &msg)
  if err != nil {
    t.Errorf("FIX ME")
  }

  digest := msg.Cid().Bytes()

  sig, err := r.Sign(digest, sk)
  if err != nil {
    t.Errorf("FIX ME")
  }

  if base64.StdEncoding.EncodeToString(sig) != "BjmEhQYMoqTeuXAn9Rj0VWk2DDhzpDA5JvppCacpnUxViDRjEgg2NY/zOWiC7g3CzxWWG9SVzfs94e4ui9N2jgE=" {
    t.Fail()
  }

}

func TestVerify(t *testing.T) {
  unsignedTx := `{
    "To": "t17uoq6tp427uzv7fztkbsnn64iwotfrristwpryy",
    "From": "t1d2xrzcslx7xlbbylc5c3d5lvandqw4iwl6epxba",
    "Nonce": 1,
    "Value": "100000",
    "GasPrice": "2500",
    "GasLimit": 25000,
    "Method": 0,
    "Params": ""
  }`

  pk, err := hex.DecodeString("0435e752dc6b4113f78edcf2cf7b8082e442021de5f00818f555397a6f181af795ace98f0f7d065793eaffa1b06bf52e572c97030c53a2396dfab40ba0e976b108")
  if err != nil {
    t.Errorf("FIX ME")
  }
  sig, err := base64.StdEncoding.DecodeString("BjmEhQYMoqTeuXAn9Rj0VWk2DDhzpDA5JvppCacpnUxViDRjEgg2NY/zOWiC7g3CzxWWG9SVzfs94e4ui9N2jgE=")
  if err != nil {
    t.Errorf("FIX ME")
  }
  r := &RosettaConstructionFilecoin{false}

  rawIn := json.RawMessage(unsignedTx)

  bytes, err := rawIn.MarshalJSON()
  if err != nil {
    t.Errorf("FIX ME")
  }

  var msg types.Message
  err = json.Unmarshal(bytes, &msg)
  if err != nil {
    t.Errorf("FIX ME")
  }

  digest := msg.Cid().Bytes()

  err = r.Verify(digest, pk, sig)

  if err != nil {
    t.Fail()
  }

}

func TestConstructPayment(t *testing.T) {
  expected := `{
    "Version": 0,
    "To": "t17uoq6tp427uzv7fztkbsnn64iwotfrristwpryy",
    "From": "t1d2xrzcslx7xlbbylc5c3d5lvandqw4iwl6epxba",
    "Nonce": 1,
    "Value": "100000",
    "GasPrice": "2500",
    "GasLimit": 25000,
    "Method": 0,
    "Params": ""
  }`
  r := &RosettaConstructionFilecoin{false}
  mtx := TxMetadata{
    Nonce: 1,
    GasPrice: "2500",
    GasLimit: 25000,
  }
  pr := &PaymentRequest{
    From: "t1d2xrzcslx7xlbbylc5c3d5lvandqw4iwl6epxba",
    To: "t17uoq6tp427uzv7fztkbsnn64iwotfrristwpryy",
    Quantity: 100000,
    Metadata: mtx,
  }

  tx, err := r.ConstructPayment(pr)
  if err != nil {
    t.Errorf("FIX ME")
  }

  var expectedTx types.Message
  var resultTx types.Message

  err = json.Unmarshal([]byte(expected), &expectedTx)
  if err != nil {
    t.Errorf("FIX ME")
  }

  err = json.Unmarshal([]byte(tx), &resultTx)
  if err != nil {
    t.Errorf("FIX ME")
  }

  if !reflect.DeepEqual(expectedTx, resultTx) {
    t.Fail()
  }

}

func TestConstructMultisigPayment(t *testing.T) {
  expected := `{
    "Version": 0,
    "To": "t01002",
    "From": "t1d2xrzcslx7xlbbylc5c3d5lvandqw4iwl6epxba",
    "Nonce": 1,
    "Value": "0",
    "GasPrice": "2500",
    "GasLimit": 25000,
    "Method": 2,
    "Params": "hFUB/R0PTfzX6Zr8uZqDJrfcRZ0yxihDAAPoAEA="
  }`
  r := &RosettaConstructionFilecoin{false}
  mtx := TxMetadata{
    Nonce: 1,
    GasPrice: "2500",
    GasLimit: 25000,
  }
  params := MultisigPaymentParams{
    To: "t17uoq6tp427uzv7fztkbsnn64iwotfrristwpryy",
    Quantity: 1000,
  }
  request := &MultisigPaymentRequest{
    Multisig: "t01002",
    From: "t1d2xrzcslx7xlbbylc5c3d5lvandqw4iwl6epxba",
    Metadata: mtx,
    Params: params,
  }

  result, err := r.ConstructMultisigPayment(request)
  if err != nil {
    t.Errorf("FIX ME")
  }

  var expectedMessage types.Message
  var resultMessage types.Message

  err = json.Unmarshal([]byte(expected), &expectedMessage)
  if err != nil {
    t.Errorf("FIX ME")
  }

  err = json.Unmarshal([]byte(result), &resultMessage)
  if err != nil {
    t.Errorf("FIX ME")
  }

  if !reflect.DeepEqual(expectedMessage, resultMessage) {
    t.Fail()
  }
}

func TestConstructSwapAuthorizedParty(t *testing.T) {
  expected := `{
    "Version": 0,
    "To": "t01002",
    "From": "t1d2xrzcslx7xlbbylc5c3d5lvandqw4iwl6epxba",
    "Nonce": 1,
    "Value": "0",
    "GasPrice": "2500",
    "GasLimit": 25000,
    "Method": 7,
    "Params": "glUB/R0PTfzX6Zr8uZqDJrfcRZ0yxihYMQOuzzY13jMOTmpShDszOIxbNhcAhlxVLRYZmVI87UlsVOZXuGJil7OSixyQSOsTXug="
  }`
  r := &RosettaConstructionFilecoin{false}
  mtx := TxMetadata{
    Nonce: 1,
    GasPrice: "2500",
    GasLimit: 25000,
  }
  params := SwapAuthorizedPartyParams{
    From: "t17uoq6tp427uzv7fztkbsnn64iwotfrristwpryy",
    To: "t3v3htmno6gmhe42ssqq5tgoemlm3boaeglrks2fqztfjdz3kjnrkomv5ymjrjpm4srmojashlcnporcluiyaa",
  }
  request := &SwapAuthorizedPartyRequest{
    Multisig: "t01002",
    From: "t1d2xrzcslx7xlbbylc5c3d5lvandqw4iwl6epxba",
    Metadata: mtx,
    Params: params,
  }

  result, err := r.ConstructSwapAuthorizedParty(request)
  if err != nil {
    t.Errorf("FIX ME")
  }

  var expectedMessage types.Message
  var resultMessage types.Message

  err = json.Unmarshal([]byte(expected), &expectedMessage)
  if err != nil {
    t.Errorf("FIX ME")
  }

  err = json.Unmarshal([]byte(result), &resultMessage)
  if err != nil {
    t.Errorf("FIX ME")
  }

  if !reflect.DeepEqual(expectedMessage, resultMessage) {
    t.Fail()
  }

}

func TestSignTx(t *testing.T) {
  unsignedTx := `{
    "To": "t17uoq6tp427uzv7fztkbsnn64iwotfrristwpryy",
    "From": "t1d2xrzcslx7xlbbylc5c3d5lvandqw4iwl6epxba",
    "Nonce": 1,
    "Value": "100000",
    "GasPrice": "2500",
    "GasLimit": 25000,
    "Method": 0,
    "Params": ""
  }`
  sk := "f15716d3b003b304b8055d9cc62e6b9c869d56cc930c3858d4d7c31f5f53f14a"
  r := &RosettaConstructionFilecoin{false}

  skBytes, err := hex.DecodeString(sk)

  if err != nil {
    t.Errorf("Invalid test case")
  }

  signedTx, err := r.SignTx(unsignedTx, skBytes)
  if err != nil {
    t.Error(err)
  }

  t.Log(signedTx)

  rawIn := json.RawMessage(signedTx)

  bytes, err := rawIn.MarshalJSON()
  if err != nil {
    t.Errorf("Not a json string")
  }

  var msg types.SignedMessage
  err = json.Unmarshal(bytes, &msg)
  if err != nil {
    t.Errorf("Not a SignedMessage")
  }

  dataSignature := base64.StdEncoding.EncodeToString(msg.Signature.Data)
  if dataSignature != "BjmEhQYMoqTeuXAn9Rj0VWk2DDhzpDA5JvppCacpnUxViDRjEgg2NY/zOWiC7g3CzxWWG9SVzfs94e4ui9N2jgE=" {
    t.Fail()
  }

}

func TestParseTx(t *testing.T) {
  serializedTx := "89005501fd1d0f4dfcd7e99afcb99a8326b7dc459d32c6285501b882619d46558f3d9e316d11b48dcf211327025a0144000186a0430009c41961a80040"
  r := &RosettaConstructionFilecoin{false}
  b, err := hex.DecodeString(serializedTx)

  if err != nil {
    t.Errorf("Invalid test case")
  }

  msg, err := r.ParseTx(b)

  if err != nil {
    t.Fail()
  }

  switch msg := msg.(type) {
    case types.Message:
      if msg.To.String() != "t17uoq6tp427uzv7fztkbsnn64iwotfrristwpryy" {
        t.Errorf("Invalid To address returned")
      }
      if msg.From.String() != "t1xcbgdhkgkwht3hrrnui3jdopeejsoas2rujnkdi" {
        t.Errorf("Invalid From address returned")
      }
      if msg.Nonce != uint64(1) {
        t.Errorf("Invalid Nonce returned")
      }
      if types.BigCmp(msg.Value, types.NewInt(100000)) > 0 {
        t.Errorf("Invalid Value returned")
      }
      if types.BigCmp(msg.GasPrice,types.NewInt(2500)) > 0 {
        t.Errorf("Invalid GasPrice returned")
      }
      if msg.GasLimit != int64(25000) {
        t.Errorf("Invalid GasLimit returned")
      }
      if msg.Method != abi.MethodNum(0) {
        t.Errorf("Invalid Method returned")
      }
      // FIXME
      /*if msg.Params != make([]byte, 0) {
        t.Errorf("Invalid Params returned")
      }*/
    case types.SignedMessage:
      t.Log(msg.Message.To)
    default:
      t.Errorf("This should never happened")

  }
}

func TestHash(t *testing.T) {
  signedTx := `{
    "Message": {
      "To": "t17uoq6tp427uzv7fztkbsnn64iwotfrristwpryy",
      "From": "t1d2xrzcslx7xlbbylc5c3d5lvandqw4iwl6epxba",
      "Nonce": 1,
      "Value": "100000",
      "GasPrice": "2500",
      "GasLimit": 25000,
      "Method": 0,
      "Params": ""
    },
    "Signature": {
      "Types": 1,
      "Data": "BjmEhQYMoqTeuXAn9Rj0VWk2DDhzpDA5JvppCacpnUxViDRjEgg2NY/zOWiC7g3CzxWWG9SVzfs94e4ui9N2jgE="
    }
  }`
  r := &RosettaConstructionFilecoin{false}

  cid, err := r.Hash(signedTx)

  if err != nil {
    t.Errorf("Something went Wrong")
  }

  t.Log(cid)

  if cid != "bafy2bzacedbhs4ewvpqjg2vdarfo4ux7nbwzvwrh36jrwsnpf6474qaicd6by" {
      t.Fail()
  }
}
