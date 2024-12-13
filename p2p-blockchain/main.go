// package main

// import (
// 	"bufio"
// 	"context"
// 	"crypto/rand"
// 	"crypto/sha256"
// 	"encoding/hex"
// 	"encoding/json"
// 	"flag"
// 	"fmt"
// 	"github.com/davecgh/go-spew/spew"
// 	golog "github.com/ipfs/go-log"
// 	gologging "github.com/ipfs/go-log/v2"
// 	libp2p "github.com/libp2p/go-libp2p"
// 	crypto "github.com/libp2p/go-libp2p/core/crypto"
// 	host "github.com/libp2p/go-libp2p/core/host"
// 	net "github.com/libp2p/go-libp2p/core/network"
// 	"github.com/libp2p/go-libp2p/core/peer"
// 	pstore "github.com/libp2p/go-libp2p/core/peerstore"
// 	ma "github.com/multiformats/go-multiaddr"
// 	"io"
// 	"log"
// 	mrand "math/rand"
// 	"os"
// 	"os/exec"
// 	"strconv"
// 	"strings"
// 	"sync"
// 	"time"
// )

// // Block represents each 'item' in the blockchain
// type Block struct {
// 	Index      int    // Block index
// 	Timestamp  string // Timestamp of block creation
// 	BPM        int    // BPM for simplicity (kept as is)
// 	Hash       string // Hash of the block
// 	PrevHash   string // Hash of the previous block
// 	Difficulty int    // Mining difficulty
// 	Nonce      string // Nonce for PoW
// 	AlgoCID    string // CID for the algorithm file (from IPFS/Kubo)
// 	DatasetCID string // CID for the dataset file (from IPFS/Kubo)
// 	Result     string // Computed result (e.g., KNN accuracy)
// }

// const difficulty = 1

// // Blockchain is a series of validated Blocks
// var Blockchain []Block
// var mutex = &sync.Mutex{}

// // makeBasicHost creates a LibP2P host with a random peer ID listening on the
// // given multiaddress. It will use secio if secio is true.
// func makeBasicHost(listenPort int, secio bool, randseed int64) (host.Host, error) {

// 	// If the seed is zero, use real cryptographic randomness. Otherwise, use a
// 	// deterministic randomness source to make generated keys stay the same
// 	// across multiple runs
// 	var r io.Reader
// 	if randseed == 0 {
// 		r = rand.Reader
// 	} else {
// 		r = mrand.New(mrand.NewSource(randseed))
// 	}

// 	// Generate a key pair for this host. We will use it
// 	// to obtain a valid host ID.
// 	priv, _, err := crypto.GenerateKeyPairWithReader(crypto.RSA, 2048, r)
// 	if err != nil {
// 		return nil, err
// 	}

// 	opts := []libp2p.Option{
// 		libp2p.ListenAddrStrings(fmt.Sprintf("/ip4/127.0.0.1/tcp/%d", listenPort)),
// 		libp2p.Identity(priv),
// 	}

// 	basicHost, err := libp2p.New(opts...)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// Build host multiaddress
// 	hostAddr, _ := ma.NewMultiaddr(fmt.Sprintf("/ipfs/%s", basicHost.ID()))

// 	// Now we can build a full multiaddress to reach this host
// 	// by encapsulating both addresses:
// 	addrs := basicHost.Addrs()
// 	var addr ma.Multiaddr
// 	// select the address starting with "ip4"
// 	for _, i := range addrs {
// 		if strings.HasPrefix(i.String(), "/ip4") {
// 			addr = i
// 			break
// 		}
// 	}
// 	fullAddr := addr.Encapsulate(hostAddr)
// 	log.Printf("I am %s\n", fullAddr)
// 	if secio {
// 		log.Printf("Now run \"go run main.go -l %d -d %s -secio\" on a different terminal\n", listenPort+1, fullAddr)
// 	} else {
// 		log.Printf("Now run \"go run main.go -l %d -d %s\" on a different terminal\n", listenPort+1, fullAddr)
// 	}

// 	return basicHost, nil
// }

// func handleStream(s net.Stream) {

// 	log.Println("Got a new stream!")

// 	// Create a buffer stream for non blocking read and write.
// 	rw := bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))

// 	go readData(rw)
// 	go writeData(rw)

// 	// stream 's' will stay open until you close it (or the other side closes it).
// }

// func readData(rw *bufio.ReadWriter) {
// 	for {
// 		str, err := rw.ReadString('\n')
// 		if err != nil {
// 			log.Fatal(err)
// 		}

// 		if str == "" {
// 			return
// 		}
// 		if str != "\n" {
// 			var block Block
// 			if err := json.Unmarshal([]byte(str), &block); err != nil {
// 				log.Println("Failed to unmarshal block:", err)
// 				continue
// 			}

// 			mutex.Lock()
// 			if isBlockValid(block) {
// 				Blockchain = append(Blockchain, block)
// 				bytes, err := json.MarshalIndent(Blockchain, "", "  ")
// 				if err != nil {
// 					log.Fatal(err)
// 				}
// 				fmt.Printf("\x1b[32m%s\x1b[0m> ", string(bytes))
// 			} else {
// 				log.Println("Block validation failed!")
// 			}
// 			mutex.Unlock()
// 		}
// 	}
// }

// // Modified writeData function to include AlgoCID and DatasetCID
// func writeData(rw *bufio.ReadWriter) {
// 	stdReader := bufio.NewReader(os.Stdin)

// 	for {
// 		fmt.Print("> Enter BPM: ")
// 		sendData, err := stdReader.ReadString('\n')
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		bpm, err := strconv.Atoi(strings.TrimSpace(sendData))
// 		if err != nil {
// 			log.Fatal(err)
// 		}

// 		fmt.Print("> Enter Algorithm CID: ")
// 		algoCID, err := stdReader.ReadString('\n')
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		fmt.Print("> Enter Dataset CID: ")
// 		datasetCID, err := stdReader.ReadString('\n')
// 		if err != nil {
// 			log.Fatal(err)
// 		}

// 		newBlock, err := generateBlock(Blockchain[len(Blockchain)-1], bpm, strings.TrimSpace(algoCID), strings.TrimSpace(datasetCID))
// 		if err != nil {
// 			log.Println("Failed to generate block:", err)
// 			continue
// 		}

// 		mutex.Lock()
// 		Blockchain = append(Blockchain, newBlock)
// 		mutex.Unlock()

// 		bytes, err := json.Marshal(newBlock)
// 		if err != nil {
// 			log.Println(err)
// 		}

// 		spew.Dump(Blockchain)

// 		mutex.Lock()
// 		rw.WriteString(fmt.Sprintf("%s\n", string(bytes)))
// 		rw.Flush()
// 		mutex.Unlock()
// 	}
// }

// func main() {
// 	t := time.Now()
// 	genesisBlock := Block{}
// 	genesisBlock = Block{0, t.String(), 0, calculateHash(genesisBlock), "", difficulty, "", "", "", ""}

// 	Blockchain = append(Blockchain, genesisBlock)

// 	// LibP2P code uses golog to log messages. They log with different
// 	// string IDs (i.e. "swarm"). We can control the verbosity level for
// 	// all loggers with:
// 	golog.SetAllLoggers(gologging.LevelInfo) // Change to DEBUG for extra info

// 	// Parse options from the command line
// 	listenF := flag.Int("l", 0, "wait for incoming connections")
// 	target := flag.String("d", "", "target peer to dial")
// 	secio := flag.Bool("secio", false, "enable secio")
// 	seed := flag.Int64("seed", 0, "set random seed for id generation")
// 	flag.Parse()

// 	if *listenF == 0 {
// 		log.Fatal("Please provide a port to bind on with -l")
// 	}

// 	// Make a host that listens on the given multiaddress
// 	ha, err := makeBasicHost(*listenF, *secio, *seed)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	if *target == "" {
// 		log.Println("listening for connections")
// 		// Set a stream handler on host A. /p2p/1.0.0 is
// 		// a user-defined protocol name.
// 		ha.SetStreamHandler("/p2p/1.0.0", handleStream)

// 		select {} // hang forever
// 		/**** This is where the listener code ends ****/
// 	} else {
// 		ha.SetStreamHandler("/p2p/1.0.0", handleStream)

// 		// The following code extracts target's peer ID from the
// 		// given multiaddress
// 		ipfsaddr, err := ma.NewMultiaddr(*target)
// 		if err != nil {
// 			log.Fatalln(err)
// 		}

// 		pid, err := ipfsaddr.ValueForProtocol(ma.P_IPFS)
// 		if err != nil {
// 			log.Fatalln(err)
// 		}

// 		peerid, err := peer.Decode(pid)
// 		if err != nil {
// 			log.Fatalln(err)
// 		}

// 		// Decapsulate the /ipfs/<peerID> part from the target
// 		// /ip4/<a.b.c.d>/ipfs/<peer> becomes /ip4/<a.b.c.d>
// 		targetPeerAddr, _ := ma.NewMultiaddr(
// 			fmt.Sprintf("/ipfs/%s", peer.ToCid(peerid)))
// 		targetAddr := ipfsaddr.Decapsulate(targetPeerAddr)

// 		// We have a peer ID and a targetAddr so we add it to the peerstore
// 		// so LibP2P knows how to contact it
// 		ha.Peerstore().AddAddr(peerid, targetAddr, pstore.PermanentAddrTTL)

// 		log.Println("opening stream")
// 		// make a new stream from host B to host A
// 		// it should be handled on host A by the handler we set above because
// 		// we use the same /p2p/1.0.0 protocol
// 		s, err := ha.NewStream(context.Background(), peerid, "/p2p/1.0.0")
// 		if err != nil {
// 			log.Fatalln(err)
// 		}
// 		// Create a buffered stream so that read and writes are non blocking.
// 		rw := bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))

// 		// Create a thread to read and write data.
// 		go writeData(rw)
// 		go readData(rw)

// 		select {} // hang forever

// 	}
// }

// // make sure block is valid by checking index, and comparing the hash of the previous block
// func calculateHash(block Block) string {
// 	record := strconv.Itoa(block.Index) + block.Timestamp + strconv.Itoa(block.BPM) + block.PrevHash + block.Nonce
// 	h := sha256.New()
// 	h.Write([]byte(record))
// 	hashed := h.Sum(nil)
// 	return hex.EncodeToString(hashed)
// }

// func isHashValid(hash string, difficulty int) bool {
// 	prefix := strings.Repeat("0", difficulty)
// 	return strings.HasPrefix(hash, prefix)
// }

// func getNonceValue(nonce string) int {
// 	if nonce == "" {
// 		return 0
// 	}

// 	value, err := strconv.Atoi(nonce)
// 	if err != nil {
// 		fmt.Println("Error converting nonce to integer:", err)
// 		return 0
// 	}

// 	return value + 1
// }

// // fetchFromIPFS fetches files from Kubo/IPFS using their CID and saves them locally
// func fetchFromIPFS(cid string) (string, error) {
// 	// Using the "ipfs cat" command to fetch the file locally
// 	cmd := exec.Command("ipfs", "cat", cid)
// 	output, err := cmd.Output()
// 	if err != nil {
// 		return "", fmt.Errorf("failed to fetch file from IPFS with CID %s: %v", cid, err)
// 	}

// 	// Save the output to a temporary file
// 	filePath := fmt.Sprintf("./%s", cid) // Save the file locally using its CID as the name
// 	err = os.WriteFile(filePath, output, 0644)
// 	if err != nil {
// 		return "", fmt.Errorf("failed to save file locally: %v", err)
// 	}
// 	return filePath, nil
// }

// // executeAlgorithm runs the algorithm on the dataset and returns the result
// func executeAlgorithm(algoPath, datasetPath string) (string, error) {
// 	// Running the algorithm using "go run" with the provided dataset
// 	cmd := exec.Command("go", "run", algoPath, datasetPath)
// 	output, err := cmd.CombinedOutput()
// 	if err != nil {
// 		return "", fmt.Errorf("failed to execute algorithm: %v, output: %s", err, string(output))
// 	}
// 	return strings.TrimSpace(string(output)), nil
// }

// func generateBlock(oldBlock Block, bpm int, algoCID, datasetCID string) (Block, error) {
// 	var newBlock Block
// 	t := time.Now()

// 	newBlock.Index = oldBlock.Index + 1
// 	newBlock.Timestamp = t.String()
// 	newBlock.BPM = bpm
// 	newBlock.PrevHash = oldBlock.Hash
// 	newBlock.Difficulty = difficulty
// 	newBlock.AlgoCID = algoCID
// 	newBlock.DatasetCID = datasetCID

// 	// Fetch the algorithm and dataset from IPFS/Kubo
// 	algoPath, err := fetchFromIPFS(algoCID)
// 	if err != nil {
// 		return newBlock, fmt.Errorf("failed to fetch algorithm from IPFS: %v", err)
// 	}
// 	datasetPath, err := fetchFromIPFS(datasetCID)
// 	if err != nil {
// 		return newBlock, fmt.Errorf("failed to fetch dataset from IPFS: %v", err)
// 	}

// 	// Execute the algorithm on the dataset
// 	result, err := executeAlgorithm(algoPath, datasetPath)
// 	if err != nil {
// 		return newBlock, fmt.Errorf("failed to execute algorithm: %v", err)
// 	}
// 	newBlock.Result = result

// 	// Perform Proof of Work
// 	for i := getNonceValue(oldBlock.Nonce); ; i++ {
// 		hex := fmt.Sprintf("%d", i)
// 		newBlock.Nonce = hex
// 		newBlock.Hash = calculateHash(newBlock)
// 		if isHashValid(newBlock.Hash, newBlock.Difficulty) {
// 			break
// 		}
// 	}

// 	return newBlock, nil
// }

// func replaceChain(newBlocks []Block) {
// 	mutex.Lock()
// 	if len(newBlocks) > len(Blockchain) {
// 		Blockchain = newBlocks
// 	}
// 	mutex.Unlock()
// }

// func isBlockValid(block Block) bool {
// 	// Fetch the algorithm and dataset from IPFS/Kubo
// 	algoPath, err := fetchFromIPFS(block.AlgoCID)
// 	if err != nil {
// 		log.Printf("Failed to fetch algorithm: %v", err)
// 		return false
// 	}
// 	datasetPath, err := fetchFromIPFS(block.DatasetCID)
// 	if err != nil {
// 		log.Printf("Failed to fetch dataset: %v", err)
// 		return false
// 	}

// 	// Execute the algorithm on the dataset
// 	result, err := executeAlgorithm(algoPath, datasetPath)
// 	if err != nil {
// 		log.Printf("Failed to execute algorithm: %v", err)
// 		return false
// 	}

// 	// Compare the computed result with the result in the block
// 	return strings.TrimSpace(result) == strings.TrimSpace(block.Result)
// }

package main

import (
	"bufio"
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	mrand "math/rand"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/davecgh/go-spew/spew"
	golog "github.com/ipfs/go-log"
	gologging "github.com/ipfs/go-log/v2"
	libp2p "github.com/libp2p/go-libp2p"
	crypto "github.com/libp2p/go-libp2p/core/crypto"
	host "github.com/libp2p/go-libp2p/core/host"
	net "github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	pstore "github.com/libp2p/go-libp2p/core/peerstore"
	ma "github.com/multiformats/go-multiaddr"
)

// Block represents each 'item' in the blockchain
type Block struct {
	Index      int
	Timestamp  string
	BPM        int
	Hash       string
	PrevHash   string
	Difficulty int
	Nounce     string
}

const difficulty = 1

// Blockchain is a series of validated Blocks
var Blockchain []Block
var mutex = &sync.Mutex{}

// makeBasicHost creates a LibP2P host with a random peer ID listening on the
// given multiaddress. It will use secio if secio is true.
func makeBasicHost(listenPort int, secio bool, randseed int64) (host.Host, error) {

	// If the seed is zero, use real cryptographic randomness. Otherwise, use a
	// deterministic randomness source to make generated keys stay the same
	// across multiple runs
	var r io.Reader
	if randseed == 0 {
		r = rand.Reader
	} else {
		r = mrand.New(mrand.NewSource(randseed))
	}

	// Generate a key pair for this host. We will use it
	// to obtain a valid host ID.
	priv, _, err := crypto.GenerateKeyPairWithReader(crypto.RSA, 2048, r)
	if err != nil {
		return nil, err
	}

	opts := []libp2p.Option{
		libp2p.ListenAddrStrings(fmt.Sprintf("/ip4/127.0.0.1/tcp/%d", listenPort)),
		libp2p.Identity(priv),
	}

	basicHost, err := libp2p.New(opts...)
	if err != nil {
		return nil, err
	}

	// Build host multiaddress
	hostAddr, _ := ma.NewMultiaddr(fmt.Sprintf("/ipfs/%s", basicHost.ID()))

	// Now we can build a full multiaddress to reach this host
	// by encapsulating both addresses:
	addrs := basicHost.Addrs()
	var addr ma.Multiaddr
	// select the address starting with "ip4"
	for _, i := range addrs {
		if strings.HasPrefix(i.String(), "/ip4") {
			addr = i
			break
		}
	}
	fullAddr := addr.Encapsulate(hostAddr)
	log.Printf("I am %s\n", fullAddr)
	if secio {
		log.Printf("Now run \"go run main.go -l %d -d %s -secio\" on a different terminal\n", listenPort+1, fullAddr)
	} else {
		log.Printf("Now run \"go run main.go -l %d -d %s\" on a different terminal\n", listenPort+1, fullAddr)
	}

	return basicHost, nil
}

func handleStream(s net.Stream) {

	log.Println("Got a new stream!")

	// Create a buffer stream for non blocking read and write.
	rw := bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))

	go readData(rw)
	go writeData(rw)

	// stream 's' will stay open until you close it (or the other side closes it).
}

func readData(rw *bufio.ReadWriter) {

	for {
		str, err := rw.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}

		if str == "" {
			return
		}
		if str != "\n" {

			chain := make([]Block, 0)
			if err := json.Unmarshal([]byte(str), &chain); err != nil {
				log.Fatal(err)
			}

			mutex.Lock()
			if len(chain) > len(Blockchain) {
				Blockchain = chain
				bytes, err := json.MarshalIndent(Blockchain, "", "  ")
				if err != nil {

					log.Fatal(err)
				}
				// Green console color: 	\x1b[32m
				// Reset console color: 	\x1b[0m
				fmt.Printf("\x1b[32m%s\x1b[0m> ", string(bytes))
			}
			mutex.Unlock()
		}
	}
}

func writeData(rw *bufio.ReadWriter) {

	go func() {
		for {
			time.Sleep(5 * time.Second)
			mutex.Lock()
			bytes, err := json.Marshal(Blockchain)
			if err != nil {
				log.Println(err)
			}
			mutex.Unlock()

			mutex.Lock()
			rw.WriteString(fmt.Sprintf("%s\n", string(bytes)))
			rw.Flush()
			mutex.Unlock()

		}
	}()

	stdReader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")
		sendData, err := stdReader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}

		// sendData = strings.Replace(sendData, "\n", "", -1)
		sendData = strings.TrimSpace(sendData)
		bpm, err := strconv.Atoi(sendData)
		if err != nil {
			log.Fatal(err)
		}
		newBlock, err := generateBlock(Blockchain[len(Blockchain)-1], bpm)

		if isBlockValid(newBlock, Blockchain[len(Blockchain)-1]) {
			mutex.Lock()
			Blockchain = append(Blockchain, newBlock)
			mutex.Unlock()
		}

		bytes, err := json.Marshal(Blockchain)
		if err != nil {
			log.Println(err)
		}

		spew.Dump(Blockchain)

		mutex.Lock()
		rw.WriteString(fmt.Sprintf("%s\n", string(bytes)))
		rw.Flush()
		mutex.Unlock()
	}

}

func main() {
	t := time.Now()
	genesisBlock := Block{}
	genesisBlock = Block{0, t.String(), 0, calculateHash(genesisBlock), "", difficulty, ""}

	Blockchain = append(Blockchain, genesisBlock)

	// LibP2P code uses golog to log messages. They log with different
	// string IDs (i.e. "swarm"). We can control the verbosity level for
	// all loggers with:
	golog.SetAllLoggers(gologging.LevelInfo) // Change to DEBUG for extra info

	// Parse options from the command line
	listenF := flag.Int("l", 0, "wait for incoming connections")
	target := flag.String("d", "", "target peer to dial")
	secio := flag.Bool("secio", false, "enable secio")
	seed := flag.Int64("seed", 0, "set random seed for id generation")
	flag.Parse()

	if *listenF == 0 {
		log.Fatal("Please provide a port to bind on with -l")
	}

	// Make a host that listens on the given multiaddress
	ha, err := makeBasicHost(*listenF, *secio, *seed)
	if err != nil {
		log.Fatal(err)
	}

	if *target == "" {
		log.Println("listening for connections")
		// Set a stream handler on host A. /p2p/1.0.0 is
		// a user-defined protocol name.
		ha.SetStreamHandler("/p2p/1.0.0", handleStream)

		select {} // hang forever
		/**** This is where the listener code ends ****/
	} else {
		ha.SetStreamHandler("/p2p/1.0.0", handleStream)

		// The following code extracts target's peer ID from the
		// given multiaddress
		ipfsaddr, err := ma.NewMultiaddr(*target)
		if err != nil {
			log.Fatalln(err)
		}

		pid, err := ipfsaddr.ValueForProtocol(ma.P_IPFS)
		if err != nil {
			log.Fatalln(err)
		}

		peerid, err := peer.Decode(pid)
		if err != nil {
			log.Fatalln(err)
		}

		// Decapsulate the /ipfs/<peerID> part from the target
		// /ip4/<a.b.c.d>/ipfs/<peer> becomes /ip4/<a.b.c.d>
		targetPeerAddr, _ := ma.NewMultiaddr(
			fmt.Sprintf("/ipfs/%s", peer.ToCid(peerid)))
		targetAddr := ipfsaddr.Decapsulate(targetPeerAddr)

		// We have a peer ID and a targetAddr so we add it to the peerstore
		// so LibP2P knows how to contact it
		ha.Peerstore().AddAddr(peerid, targetAddr, pstore.PermanentAddrTTL)

		log.Println("opening stream")
		// make a new stream from host B to host A
		// it should be handled on host A by the handler we set above because
		// we use the same /p2p/1.0.0 protocol
		s, err := ha.NewStream(context.Background(), peerid, "/p2p/1.0.0")
		if err != nil {
			log.Fatalln(err)
		}
		// Create a buffered stream so that read and writes are non blocking.
		rw := bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))

		// Create a thread to read and write data.
		go writeData(rw)
		go readData(rw)

		select {} // hang forever

	}
}

// make sure block is valid by checking index, and comparing the hash of the previous block
func calculateHash(block Block) string {
	record := strconv.Itoa(block.Index) + block.Timestamp + strconv.Itoa(block.BPM) + block.PrevHash + block.Nounce
	h := sha256.New()
	h.Write([]byte(record))
	hashed := h.Sum(nil)
	return hex.EncodeToString(hashed)
}

func isHashValid(hash string, difficulty int) bool {
	prefix := strings.Repeat("0", difficulty)
	return strings.HasPrefix(hash, prefix)
}

func getNonceValue(nonce string) int {
	if nonce == "" {
		return 0
	}

	value, err := strconv.Atoi(nonce)
	if err != nil {
		fmt.Println("Error converting nonce to integer:", err)
		return 0
	}

	return value + 1
}

func generateBlock(oldBlock Block, BPM int) (Block, error) {
	var newBlock Block

	t := time.Now()

	newBlock.Index = oldBlock.Index + 1
	newBlock.Timestamp = t.String()
	newBlock.BPM = BPM
	newBlock.PrevHash = oldBlock.Hash
	newBlock.Difficulty = difficulty

	for i := getNonceValue(oldBlock.Nounce); ; i++ {
		hex := fmt.Sprintf("%d", i)
		newBlock.Nounce = hex
		if !isHashValid(calculateHash(newBlock), newBlock.Difficulty) {
			fmt.Println(calculateHash(newBlock), " do more work!!!!")
			time.Sleep(time.Second)
			continue
		} else {
			fmt.Println(calculateHash(newBlock), " work done...!!!!")
			newBlock.Hash = calculateHash(newBlock)
			break
		}
	}

	return newBlock, nil
}

func replaceChain(newBlocks []Block) {
	mutex.Lock()
	if len(newBlocks) > len(Blockchain) {
		Blockchain = newBlocks
	}
	mutex.Unlock()
}

func isBlockValid(newBlock, oldBlock Block) bool {
	if oldBlock.Index+1 != newBlock.Index {
		return false
	}

	if oldBlock.Hash != newBlock.PrevHash {
		return false
	}

	if calculateHash(newBlock) != newBlock.Hash {
		return false
	}

	return true
}
