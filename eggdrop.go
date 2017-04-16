package main

import(
  "fmt"
  "net/http"
  "encoding/json"
  "io/ioutil"
  "strconv"
  "bytes"
  "math"
)

type eggDrop struct {
  BuildingHeight int
  EggHeight int
}

func main() {

  var buildingHeight, heightToSet, maxUnbrokenHeight int
  again := "y"

  for again == "y" {

    buildingHeight = getBuildingHeight()

    fmt.Println("\nEnter building height:")
    _, err := fmt.Scanf("%d", &heightToSet)
    if err != nil {
      fmt.Println("error:", err)
    }

    if buildingHeight != heightToSet {
      buildingHeight = setBuildingHeight(heightToSet)
    }

    maxUnbrokenHeight = findMaxUnbrokenHeight(buildingHeight)

    fmt.Println("\nFor current building of height:                  ", buildingHeight)
    fmt.Println("Max height egg can be dropped without breaking:   ", maxUnbrokenHeight)

    fmt.Println("\nAgain? [y/n]:")
    _, err = fmt.Scan(&again)
    if err != nil {
      fmt.Println("error:", err)
    }
  }
}

var client = &http.Client{}

func setBuildingHeight(h int) int {
  url := "http://codetest.matter6.com/eggs/height"

  var building map[string]string
  building = make(map[string]string)

  height := strconv.Itoa(h)
  building["height"] = height

  b, err := json.Marshal(building)
  if err != nil {
    fmt.Println("error:", err)
  }

  resp, err := client.Post(url, "application/json", bytes.NewBuffer(b))
  if err != nil {
    fmt.Println("error:", err)
  }
  defer resp.Body.Close()

  return h
}

func getBuildingHeight() int {
  url := "http://codetest.matter6.com/eggs/height"

  resp, err := client.Get(url)
  if err != nil {
    fmt.Println("error:", err)
  }
  defer resp.Body.Close()

  body, err := ioutil.ReadAll(resp.Body)
  if err != nil {
    fmt.Println("error:", err)
  }

  var building map[string]string
  building = make(map[string]string)
  err = json.Unmarshal(body, &building)
  if err != nil {
    panic(error.Error)
  }

  h, _ := strconv.Atoi(building["height"])

  return h
}

func dropEgg(h int) int {
  url := "http://codetest.matter6.com/eggs/drop"

  var drop map[string]string
  drop = make(map[string]string)

  height := strconv.Itoa(h)
  drop["eggHeight"] = height

  d, err := json.Marshal(drop)
  if err != nil {
    fmt.Println("error:", err)
  }

  resp, err := client.Post(url, "application/json", bytes.NewBuffer(d))
  if err != nil {
    fmt.Println("error:", err)
  }
  defer resp.Body.Close()

  body, err := ioutil.ReadAll(resp.Body)
  if err != nil {
    fmt.Println("error:", err)
  }

  type Result struct {
    EggHeight string
    Result int
  }
  var result Result

  err = json.Unmarshal(body, &result)
  if err != nil {
    panic(error.Error)
  }

  fmt.Println("egg status:", result.Result)
  return result.Result
}

func findMaxUnbrokenHeight(buildingHeight int) int {
  var dropHeightDelta int
  var dropHeight int
  dropMax := buildingHeight
  dropMin := 1
  eggStatus := 0

  setDropHeight := func() {
    oldDropHeight := dropHeight
    dropHeight = ((dropMax - dropMin) / 2) + dropMin
    dropHeightDelta = int(math.Abs(float64(oldDropHeight - dropHeight)))
  }
  setDropHeight()

  fmt.Println("\nDropping Eggs...\n")

  for dropHeightDelta > 0 {
    fmt.Println("dropHeight:", dropHeight)

    eggStatus = dropEgg(dropHeight)
    if eggStatus == 0 {
      dropMin = dropHeight
      setDropHeight()
    } else {
      dropMax = dropHeight
      setDropHeight()
    }
  }

  return dropHeight
}
