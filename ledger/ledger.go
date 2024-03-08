package ledger

//import ("encoding/json", log, fmt)

type ledger struct {
	BackupMaster []string `json:"backupMaster"` //maby int instead og string if use of IP?
	Alive        []string `json:"alive"`
	Orders       []string `json:"orders"`
}

// Exampel of serializing (Marshal):
// type Person struct {
//     Navn  string `json:"navn"`
//     Alder int    `json:"alder"`
// }

// func main() {
//     p := Person{"Ola Nordmann", 30}
//     pJson, err := json.Marshal(p)
//     if err != nil {
//         fmt.Println(err)
//     }
//     fmt.Println(string(pJson)) // Output: {"navn":"Ola Nordmann","alder":30}
// }
//
