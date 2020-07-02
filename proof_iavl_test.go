package iavl

import (
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	db "github.com/tendermint/tm-db"
)

func TestProofOp(t *testing.T) {
	tree, err := NewMutableTreeWithOpts(db.NewMemDB(), 0, nil)
	require.NoError(t, err)
	keys := []byte{0x0a, 0x11, 0x2e, 0x32, 0x50, 0x72, 0x99, 0xa1, 0xe4, 0xf7} // 10 total.
	for _, ikey := range keys {
		key := []byte{ikey}
		tree.Set(key, key)
	}
	root := tree.WorkingHash()

	testcases := []struct {
		key           byte
		expectPresent bool
		expectProofOp string
	}{
		{0x00, false, "aa010aa7010a280808100a18012a2022b4e34a1778d6a03aac39f00d89deb886e0cc37454e300b7aebeb4f4939c0790a280804100418012a20734fad809673ab2b9672453a8b2bc8c9591e2d1d97933df5b4c3b0531bf82e720a280802100218012a20154b101a72acffe0f5e65d1e144a57dc6f97758d2049821231f02b6a5b44fe811a270a010a122001ba4719c80b6fe911b091a7c05124b64eeece964e09c058ef8f9805daca546b1801"},
		{0x0a, true, "aa010aa7010a280808100a18012a2022b4e34a1778d6a03aac39f00d89deb886e0cc37454e300b7aebeb4f4939c0790a280804100418012a20734fad809673ab2b9672453a8b2bc8c9591e2d1d97933df5b4c3b0531bf82e720a280802100218012a20154b101a72acffe0f5e65d1e144a57dc6f97758d2049821231f02b6a5b44fe811a270a010a122001ba4719c80b6fe911b091a7c05124b64eeece964e09c058ef8f9805daca546b1801"},
		{0x0b, false, "d5010ad2010a280808100a18012a2022b4e34a1778d6a03aac39f00d89deb886e0cc37454e300b7aebeb4f4939c0790a280804100418012a20734fad809673ab2b9672453a8b2bc8c9591e2d1d97933df5b4c3b0531bf82e720a280802100218012a20154b101a72acffe0f5e65d1e144a57dc6f97758d2049821231f02b6a5b44fe8112001a270a010a122001ba4719c80b6fe911b091a7c05124b64eeece964e09c058ef8f9805daca546b18011a270a011112204a64a107f0cb32536e5bce6c98c393db21cca7f4ea187ba8c4dca8b51d4ea80a1801"},
		{0x11, true, "aa010aa7010a280808100a18012a2022b4e34a1778d6a03aac39f00d89deb886e0cc37454e300b7aebeb4f4939c0790a280804100418012a20734fad809673ab2b9672453a8b2bc8c9591e2d1d97933df5b4c3b0531bf82e720a28080210021801222053d2828f35e33aecab8e411a40afb0475288973b96aed2220e9894f43a5375ad1a270a011112204a64a107f0cb32536e5bce6c98c393db21cca7f4ea187ba8c4dca8b51d4ea80a1801"},
		{0x60, false, "d5010ad2010a280808100a18012220e39776faa9ef2b83ae828860d24f807efab321d02b78081c0e68e1bf801b0e220a280806100618012a20631b10ce49ece4cc9130befac927865742fb11caf2e8fc08fc00a4a25e4bc7940a280802100218012a207a4a97f565ae0b3ea8abf175208f176ac8301665ac2d26c89be3664f90e23da612001a270a015012205c62e091b8c0565f1bafad0dad5934276143ae2ccef7a5381e8ada5b1a8d26d218011a270a01721220454349e422f05297191ead13e21d3db520e5abef52055e4964b82fb213f593a11801"},
		{0x72, true, "aa010aa7010a280808100a18012220e39776faa9ef2b83ae828860d24f807efab321d02b78081c0e68e1bf801b0e220a280806100618012a20631b10ce49ece4cc9130befac927865742fb11caf2e8fc08fc00a4a25e4bc7940a28080210021801222035f8ea805390e084854f399b42ccdeaea33a1dedc115638ac48d0600637dba1f1a270a01721220454349e422f05297191ead13e21d3db520e5abef52055e4964b82fb213f593a11801"},
		{0x99, true, "d4010ad1010a280808100a18012220e39776faa9ef2b83ae828860d24f807efab321d02b78081c0e68e1bf801b0e220a2808061006180122201d6b29f2c439fc9f15703eb7031e4a216002ea36ee9496583f97b20302b6a74e0a280804100418012a2043b83a6acefd4fd33970d1bc8fc47bed81220c752b8de7053e8ee082a2c7c1290a280802100218012a208f69a1db006c0ee9fad3c7c624b92acc88e9ed00771976ea24a64796c236fef01a270a01991220fd9528b920d6d3956e9e16114523e1889c751e8c1e040182116d4c906b43f5581801"},
		{0xaa, false, "a9020aa6020a280808100a18012220e39776faa9ef2b83ae828860d24f807efab321d02b78081c0e68e1bf801b0e220a2808061006180122201d6b29f2c439fc9f15703eb7031e4a216002ea36ee9496583f97b20302b6a74e0a280804100418012a2043b83a6acefd4fd33970d1bc8fc47bed81220c752b8de7053e8ee082a2c7c1290a280802100218012220a303930ca8831618ac7e4ddd10546cfc366fb730d6630c030a97226bbefc6935122a0a280802100218012a2077ad141b2010cf7107de941aac5b46f44fa4f41251076656a72308263a964fb91a270a01a112208a8950f7623663222542c9469c73be3c4c81bbdf019e2c577590a61f2ce9a15718011a270a01e412205e1effe9b7bab73dce628ccd9f0cbbb16c1e6efc6c4f311e59992a467bc119fd1801"},
		{0xe4, true, "d4010ad1010a280808100a18012220e39776faa9ef2b83ae828860d24f807efab321d02b78081c0e68e1bf801b0e220a2808061006180122201d6b29f2c439fc9f15703eb7031e4a216002ea36ee9496583f97b20302b6a74e0a2808041004180122208bc4764843fdd745dc853fa62f2fac0001feae9e46136192f466c09773e2ed050a280802100218012a2077ad141b2010cf7107de941aac5b46f44fa4f41251076656a72308263a964fb91a270a01e412205e1effe9b7bab73dce628ccd9f0cbbb16c1e6efc6c4f311e59992a467bc119fd1801"},
		{0xf7, true, "d4010ad1010a280808100a18012220e39776faa9ef2b83ae828860d24f807efab321d02b78081c0e68e1bf801b0e220a2808061006180122201d6b29f2c439fc9f15703eb7031e4a216002ea36ee9496583f97b20302b6a74e0a2808041004180122208bc4764843fdd745dc853fa62f2fac0001feae9e46136192f466c09773e2ed050a28080210021801222032af6e3eec2b63d5fe1bd992a89ef3467b3cee639c068cace942f01326098f171a270a01f7122050868f20258bbc9cce0da2719e8654c108733dd2f663b8737c574ec0ead93eb31801"},
		{0xfe, false, "d4010ad1010a280808100a18012220e39776faa9ef2b83ae828860d24f807efab321d02b78081c0e68e1bf801b0e220a2808061006180122201d6b29f2c439fc9f15703eb7031e4a216002ea36ee9496583f97b20302b6a74e0a2808041004180122208bc4764843fdd745dc853fa62f2fac0001feae9e46136192f466c09773e2ed050a28080210021801222032af6e3eec2b63d5fe1bd992a89ef3467b3cee639c068cace942f01326098f171a270a01f7122050868f20258bbc9cce0da2719e8654c108733dd2f663b8737c574ec0ead93eb31801"},
		//{0xff, false, ""}, // FIXME This panics, see https://github.com/cosmos/iavl/issues/286
	}

	for _, tc := range testcases {
		tc := tc
		t.Run(fmt.Sprintf("%02x", tc.key), func(t *testing.T) {
			key := []byte{tc.key}
			value, proof, err := tree.GetWithProof(key)
			require.NoError(t, err)

			// Verify that proof is valid.
			err = proof.Verify(root)
			require.NoError(t, err)

			// Encode and decode proof, either ValueOp or AbsentOp depending on key existence.
			expectBytes, err := hex.DecodeString(tc.expectProofOp)
			require.NoError(t, err)

			if tc.expectPresent {
				require.NotNil(t, value)
				err = proof.VerifyItem(key, value)
				require.NoError(t, err)

				valueOp := NewValueOp(key, proof)
				proofOp := valueOp.ProofOp()
				assert.Equal(t, ProofOp{
					Type: ProofOpIAVLValue,
					Key:  key,
					Data: expectBytes,
				}, proofOp)

				//t.Logf("Expect: %x", expectBytes)
				//t.Logf("Actual: %x", proofOp.Data)

				d, e := ValueOpDecoder(proofOp)
				require.NoError(t, e)
				decoded := d.(ValueOp)
				err = decoded.Proof.Verify(root)
				require.NoError(t, err)
				assert.Equal(t, valueOp, decoded)

			} else {
				require.Nil(t, value)
				err = proof.VerifyAbsence(key)
				require.NoError(t, err)

				absenceOp := NewAbsenceOp(key, proof)
				proofOp := absenceOp.ProofOp()
				assert.Equal(t, ProofOp{
					Type: ProofOpIAVLAbsence,
					Key:  key,
					Data: expectBytes,
				}, proofOp)

				d, e := AbsenceOpDecoder(proofOp)
				require.NoError(t, e)
				decoded := d.(AbsenceOp)
				err = decoded.Proof.Verify(root)
				require.NoError(t, err)
				assert.Equal(t, absenceOp, decoded)
			}
		})
	}
}