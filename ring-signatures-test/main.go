// Package main allows you to generate and verify ring signatures.
package main

import (
	crand "crypto/rand"
	"fmt"
	"os"
        "runtime"
        "time"

	"github.com/t-bast/ring-signatures/ring"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.EnableBashCompletion = true
	app.Name = "ring-signatures"
	app.Usage = "generate and verify ring signatures."
	app.Version = "0.1.0"

	app.Commands = []cli.Command{
		{
			Name:      "generate",
			Aliases:   []string{"g"},
			Usage:     "generate test sequence",
			UsageText: "ring-signatures generate",
			Action:    generate,
                        Flags: []cli.Flag{
                                cli.StringFlag{
                                        Name:  "message, m",
                                        Usage: "message to sign or verify",
                                },
                                cli.IntFlag{
                                        Name:  "decoy, d",
                                        Usage: "Number of decoys",
                                },
                        },

		},
		{
			Name:    "sign",
			Aliases: []string{"s"},
			Usage:   "sign a message with a ring",
			UsageText: "Alice has private key \"Pr1v4T3k3y\", public key \"4l1c3\" and wants to sign the message \"hello!\".\n" +
				"   She wants to use Bob and Carol's public keys to form a ring.\n" +
				"   Bob's public key is \"b0b\" and Carol's public key is \"c4r0l\".\n" +
				"   Alice can form the ring [c4r0l, 4l1c3, b0b] and hide herself in that ring with the following command:\n" +
				"   ring-signatures sign --message \"hello!\" --private-key 4l1c3" +
				" --ring-index 1 --ring c4r0l --ring 4l1c3 --ring b0b",
			Action: sign,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "message, m",
					Usage: "message to sign or verify",
				},
				cli.StringFlag{
					Name:  "private-key, k",
					Usage: "private key to use for signing",
				},
				cli.IntFlag{
					Name:  "ring-index, i",
					Usage: "index of your private key in the signing ring",
				},
				cli.StringSliceFlag{
					Name:  "ring, r",
					Usage: "comma-separated list of public keys to use as ring",
				},
			},
		},
		{
			Name:      "verify",
			Aliases:   []string{"v"},
			Usage:     "verify a message signature",
			UsageText: "ring-signatures verify --message \"hello!\" --signature s1GN4tUr3",
			Action:    verify,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "message, m",
					Usage: "message to sign or verify",
				},
				cli.StringFlag{
					Name:  "signature, s",
					Usage: "signature to verify",
				},
			},
		},
	}

	app.Run(os.Args)
}

func PrintMemUsage() {
  var m runtime.MemStats
  runtime.ReadMemStats(&m)
  fmt.Printf("Alloc = %v KB", m.Alloc / 1024)
  fmt.Printf("\tTotalAlloc = %v KB\n", m.TotalAlloc / 1024)
}

func generate(c *cli.Context) error {
        decoy := c.Int("decoy")
        m := c.String("message")
        
	fmt.Println("Generating your public and private keys...")

        var pks []ring.PublicKey
        var sks []ring.PrivateKey
        for i := 0; i < decoy; i++ { 
	  pk, sk := ring.Generate(crand.Reader)
	  pks = append(pks, pk)
          sks = append(sks, sk)
        }
        
        fmt.Println("Signing message...")
        start := time.Now()
	sig, err := sks[0].Sign(crand.Reader, []byte(m), pks, 0)
	if err != nil {
		return cli.NewExitError(err, 1)
	}

	sigStr, err := sig.Encode()
	if err != nil {
		return cli.NewExitError(err, 1)
	}
        elapsed := time.Since(start) 
	fmt.Println(sigStr)
        fmt.Printf("time taken %dms\n", elapsed.Nanoseconds()/1000000)
        PrintMemUsage()
        
        fmt.Println("Verifying message...")
        start = time.Now()
        sig = &ring.Signature{}
	err = sig.Decode(sigStr)
	if err != nil {
          return cli.NewExitError("invalid signature", 1)
	}

	valid := sig.Verify([]byte(m))
	if !valid {
	  return cli.NewExitError("invalid signature", 1)
	}
        elapsed = time.Since(start)
        fmt.Println("Signature is valid.")
        fmt.Printf("time taken %dms\n", elapsed.Nanoseconds()/1000)
        PrintMemUsage()
 
	return nil
}

func sign(c *cli.Context) error {
	r := c.StringSlice("ring")
	if len(r) == 0 {
		return cli.NewExitError("you need to specify a ring to use for signing", 1)
	}

	var ringKeys []ring.PublicKey
	for _, key := range r {
		pkBytes, err := ring.ConfigDecodeKey(key)
		if err != nil {
			return cli.NewExitError(fmt.Sprintf("invalid public key: %s", key), 1)
		}

		ringKeys = append(ringKeys, ring.PublicKey(pkBytes))
	}

	m := c.String("message")
	if len(m) == 0 {
		return cli.NewExitError("you need to specify a message to sign", 1)
	}

	i := c.Int("ring-index")
	if i < 0 {
		return cli.NewExitError("invalid index", 1)
	}

	pk := c.String("private-key")
	if len(pk) == 0 {
		return cli.NewExitError("you need to specify the private key to use for signing", 1)
	}

	privKeyBytes, err := ring.ConfigDecodeKey(pk)
	if err != nil {
		return cli.NewExitError("invalid private key", 1)
	}

	privKey := ring.PrivateKey(privKeyBytes)

	fmt.Println("Signing message...")
	sig, err := privKey.Sign(crand.Reader, []byte(m), ringKeys, i)
	if err != nil {
		return cli.NewExitError(err, 1)
	}

	sigStr, err := sig.Encode()
	if err != nil {
		return cli.NewExitError(err, 1)
	}

	fmt.Println(sigStr)

	return nil
}

func verify(c *cli.Context) error {
	sigStr := c.String("signature")
	if len(sigStr) == 0 {
		return cli.NewExitError("you need to specify the signature to verify", 1)
	}

	m := c.String("message")
	if len(m) == 0 {
		return cli.NewExitError("you need to specify the signed message", 1)
	}

	sig := &ring.Signature{}
	err := sig.Decode(sigStr)
	if err != nil {
		return cli.NewExitError("invalid signature", 1)
	}

	valid := sig.Verify([]byte(m))
	if !valid {
		return cli.NewExitError("invalid signature", 1)
	}

	fmt.Println("Signature is valid.")

	return nil
}
