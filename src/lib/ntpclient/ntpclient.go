package ntpclient

import (
	_ "embed"
	"encoding/json"
	"errors"
	"io"
	"net"
	"runtime"
	"time"

	"tinygo.org/x/drivers/netlink"
	"tinygo.org/x/drivers/netlink/probe"
)

// 👁️ https://github.com/tinygo-org/drivers/blob/release/examples/net/ntpclient/main.go

//go:embed secrets.json
var secrets []byte

type Secrets struct {
	SSID     string `json:"ssid"`
	Password string `json:"pass"`
}

const NTP_PACKET_SIZE = 48

var response = make([]byte, NTP_PACKET_SIZE)

// 🟧 Synchronize system time with NTP

func SyncSystemTime(ntpHost string) {

	var s Secrets
	if err := json.Unmarshal(secrets, &s); err != nil {
		panic(err)
	}

	println("🐞 attempting to connect to SSID", s.SSID)

	link, _ := probe.Probe()

	params := netlink.ConnectParams{
		Ssid:       s.SSID,
		Passphrase: s.Password,
	}

	if err := link.NetConnect(&params); err != nil {
		println("🔥 failed to connect", err.Error())
		panic(err)
	}

	conn, err := net.Dial("udp", ntpHost)
	if err != nil {
		println("🔥 UDP dial failed", err.Error())
		panic(err)
	}

	println("🐞 requesting NTP time")

	t, err := getCurrentTime(conn)
	if err != nil {
		println("🐞 failed to get current time", err.Error())
	}

	conn.Close()
	link.NetDisconnect()

	runtime.AdjustTimeOffset(-1 * int64(time.Since(t)))

	println("🐞 system time synchronized", time.Now().String())
}

// 🟦 Helpers

func getCurrentTime(conn net.Conn) (time.Time, error) {
	if err := sendNTPpacket(conn); err != nil {
		return time.Time{}, err
	}
	n, err := conn.Read(response)
	if err != nil && err != io.EOF {
		return time.Time{}, err
	}
	if n != NTP_PACKET_SIZE {
		return time.Time{}, errors.New("expected NTP packet size")
	}
	return parseNTPpacket(response), nil
}

func sendNTPpacket(conn net.Conn) error {
	var request = [48]byte{0xe3}
	_, err := conn.Write(request[:])
	return err
}

func parseNTPpacket(r []byte) time.Time {
	// 👇 the timestamp starts at byte 40 of the received packet
	//    and is four bytes, this is NTP time (seconds since Jan 1 1900)
	t := uint32(r[40])<<24 | uint32(r[41])<<16 | uint32(r[42])<<8 | uint32(r[43])
	const seventyYears = 2208988800
	return time.Unix(int64(t-seventyYears), 0)
}
