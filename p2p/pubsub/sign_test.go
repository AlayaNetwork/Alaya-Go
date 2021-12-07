package pubsub

//import (
//	"crypto/ecdsa"
//	"github.com/AlayaNetwork/Alaya-Go/p2p/enode"
//	"github.com/AlayaNetwork/Alaya-Go/p2p/pubsub/message"
//	"testing"
//
//	pb "github.com/libp2p/go-libp2p-pubsub/pb"
//
//	"github.com/libp2p/go-libp2p-core/crypto"
//	"github.com/libp2p/go-libp2p-core/peer"
//)
//
//func TestSigning(t *testing.T) {
//	privk, _, err := crypto.GenerateKeyPair(crypto.RSA, 2048)
//	if err != nil {
//		t.Fatal(err)
//	}
//	testSignVerify(t, privk)
//
//	privk, _, err = crypto.GenerateKeyPair(crypto.Ed25519, 0)
//	if err != nil {
//		t.Fatal(err)
//	}
//	testSignVerify(t, privk)
//}
//
//func testSignVerify(t *testing.T, privk *ecdsa.PublicKey) {
//	id := enode.PubkeyToIDV4(privk)
//	topic := "foo"
//	m := message.Message{
//		Data:  []byte("abc"),
//		Topic: &topic,
//		From:  id,
//		Seqno: []byte("123"),
//	}
//	signMessage(id, privk, &m)
//	err = verifyMessageSignature(&m)
//	if err != nil {
//		t.Fatal(err)
//	}
//}
