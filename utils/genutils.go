package utils

import (
	"bytes"
	jsoniter "github.com/json-iterator/go"
	"math/rand"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/oklog/ulid"
)

// RandomLife ...
type RandomLife struct {
	SeededRand *rand.Rand
}

// Letters of the alphabet in upper and lower case
const (
	ALPHABET = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	DIGITS   = "0123456789"
)

// NewRnd ...
func NewRnd() RandomLife {
	return RandomLife{
		SeededRand: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

//GenUlid - Generates a ulid
func (rnd *RandomLife) GenUlid() string {
	return rnd.genUlid()
}

//GenerateString - Generates a string of n characters
func (rnd *RandomLife) GenerateString(n int) string {
	return rnd.generateString(n)
}

//GenerateSentence - Generates a sentence of words.
func (rnd *RandomLife) GenerateSentence(numWords int, maxWordLen int) string {

	var buffy bytes.Buffer
	var mu sync.Mutex

	// lock/unlock when accessing the rand from a goroutine
	mu.Lock()
	minLen := 2

	for count := 0; count < numWords; count++ {
		wordLen := minLen + rnd.SeededRand.Intn(maxWordLen-minLen)
		rnd.generateString(wordLen)
		buffy.WriteString(rnd.generateString(wordLen))
		buffy.WriteString(" ")
	}

	name := buffy.String()
	mu.Unlock()

	return name

}

//GenerateEmail - Generates a sentence of words.
func (rnd *RandomLife) GenerateEmail(numWords int, maxWordLen int) string {

	//var providers []string = []string{"yahoo.com", "gmail.com", "consultant.com", "googlemail.com", "hotmail.com"}

	providers := []string{"yahoo.com", "gmail.com", "consultant.com", "googlemail.com", "hotmail.com"}

	rndLen := 5 + rnd.NextInt(7)
	userName := rnd.generateString(rndLen)

	email := userName + "@" + providers[rnd.NextInt(len(providers))]

	return email

}

//NextInt - Generates a number between 0 and max, max. excluded
func (rnd *RandomLife) NextInt(max int) int {

	var mu sync.Mutex

	// lock/unlock when accessing the rand from a goroutine
	mu.Lock()
	i := 0 + rnd.SeededRand.Intn(max)
	mu.Unlock()

	return i

}

//NextFloat - Generates a number between 0 and 1
func (rnd *RandomLife) NextFloat() float64 {

	var mu sync.Mutex

	// lock/unlock when accessing the rand from a goroutine
	mu.Lock()
	i := 0 + rnd.SeededRand.Float64()
	mu.Unlock()

	return i

}

//GenerateFullName - Generates first name and the last name.
func (rnd *RandomLife) GenerateFullName(minLen int) string {

	var mu sync.Mutex

	// lock/unlock when accessing the rand from a goroutine
	mu.Lock()
	i := minLen + rnd.SeededRand.Intn(minLen/2)
	j := minLen + rnd.SeededRand.Intn(minLen/2)

	name := rnd.generateString(i) + " " + rnd.generateString(j)
	mu.Unlock()

	return name

}

//GenULID -
func (rnd *RandomLife) GenULID() string {
	return rnd.genUlid()
}

//GenTin -
func (rnd *RandomLife) GenTin() string {
	//10000028-5106
	//00000004-0005

	var mu sync.Mutex

	// lock/unlock when accessing the rand from a goroutine
	mu.Lock()
	firstPart := "000000" + strconv.Itoa(10+rnd.SeededRand.Intn(90))

	genIntForSecondPart := strconv.Itoa(1 + rnd.SeededRand.Intn(9999))

	for len(genIntForSecondPart) < 4 {
		genIntForSecondPart = "0" + genIntForSecondPart
	}

	tin := firstPart + "-" + genIntForSecondPart

	mu.Unlock()

	return tin

}

//GetArrEntryRndInt -
func (rnd *RandomLife) GetArrEntryRndInt(arr []int) int {

	var mu sync.Mutex

	// lock/unlock when accessing the rand from a goroutine
	mu.Lock()
	i := rnd.SeededRand.Intn(len(arr))

	name := arr[i]
	mu.Unlock()

	return name
}

//GetArrEntryRnd -
func (rnd *RandomLife) GetArrEntryRnd(arr []string) string {

	var mu sync.Mutex

	// lock/unlock when accessing the rand from a goroutine
	mu.Lock()
	i := rnd.SeededRand.Intn(len(arr))

	name := arr[i]
	mu.Unlock()

	return name
}

func (rnd *RandomLife) genUlid() string {
	t := time.Now().UTC()

	var mu sync.Mutex

	// lock/unlock when accessing the rand from a goroutine
	mu.Lock()
	entropy := rand.New(rand.NewSource(t.UnixNano()))
	mu.Unlock()

	id := ulid.MustNew(ulid.Timestamp(t), entropy)

	return id.String()
}

func (rnd *RandomLife) GenerateAlphaNumericString(chars int) string {
	var buff bytes.Buffer

	slice := make([]string, 0)

	slice = strings.Split(ALPHABET+DIGITS, "")

	length := len(slice)

	var mu sync.Mutex
	for ctr := 0; ctr < chars; ctr++ {

		// lock/unlock when accessing the rand from a goroutine
		mu.Lock()
		lent := rnd.SeededRand.Intn(length)
		mu.Unlock()
		buff.WriteString(string(slice[lent]))
	}

	return buff.String()

}

func (rnd *RandomLife) generateString(chars int) string {
	var buff bytes.Buffer
	length := len(ALPHABET)
	var mu sync.Mutex
	for ctr := 0; ctr < chars; ctr++ {

		// lock/unlock when accessing the rand from a goroutine
		mu.Lock()
		lent := rnd.SeededRand.Intn(length)
		mu.Unlock()
		buff.WriteString(string(ALPHABET[lent]))
	}

	return buff.String()

}

//GeneratePhone -
func (rnd *RandomLife) GeneratePhone(chars int) string {
	var buff bytes.Buffer
	length := len(DIGITS)
	var mu sync.Mutex

	arr := []string{"070", "080", "090", "081"}

	buff.WriteString(arr[rnd.SeededRand.Intn(len(arr))])

	for ctr := 0; ctr < chars-3; ctr++ {

		// lock/unlock when accessing the rand from a goroutine
		mu.Lock()
		lent := rnd.SeededRand.Intn(length)
		mu.Unlock()
		buff.WriteString(string(DIGITS[lent]))
	}

	return buff.String()

}

/**
 * startwith is a sequence of digits that the string will
 * start with.
 * chars ...The total number of characters in the
 * generated string of digits
 *
 */
func (rnd *RandomLife) GenerateDigits(startWith string, chars int) string {
	var buff bytes.Buffer
	length := len(DIGITS)
	var mu sync.Mutex

	buff.WriteString(startWith)

	for ctr := 0; ctr < chars-len(startWith); ctr++ {

		// lock/unlock when accessing the rand from a goroutine
		mu.Lock()
		lent := rnd.SeededRand.Intn(length)
		mu.Unlock()
		buff.WriteString(string(DIGITS[lent]))
	}

	return buff.String()

}

// GenerateBool ...
func (rnd *RandomLife) GenerateBool() bool {
	return rnd.SeededRand.Intn(2) == 1
}

// GenerateRandomTimestampSince ...
func (rnd *RandomLife) GenerateRandomTimestampSince(sinceYears int) int {

	yearsInSecs := int64(1000 * sinceYears * 365 * 86400)

	args := int64(time.Now().Unix()*1000 - yearsInSecs)

	if args < 0 {
		args = -1 * args
	}
	randomTime := rnd.SeededRand.Int63n(args) + yearsInSecs

	return int(randomTime)
}

// GenerateRndFloat ...Supply min and max
func (rnd *RandomLife) GenerateRndFloat(min float32, max float32) float32 {
	return min + rnd.SeededRand.Float32()*(max-min)
}

// CurrentTimeStamp  The time now
func CurrentTimeStamp() int {
	return int(time.Now().UnixNano() / 1000000)
}

func GetTimeFromTimeStamp(timestamp int) (time.Time, error) {
	i, err := strconv.ParseInt(strconv.Itoa(timestamp/1000), 10, 64)
	if err != nil {
		return time.Time{}, err
	}
	tm := time.Unix(i, 0)

	return tm, nil
}

func GetTimeStampFromRfc3339String(rfc3339DateStr string) (int, error) {

	jsonify :=
		`{
"the_time": "` + rfc3339DateStr + `"` +
			`}`

	//2017-05-24T05:56:12.000Z
	type Date struct {
		RfcDateTime time.Time `json:"the_time, string"` //t.Format(time.RFC3339)
	}

	var dt Date

	err := jsoniter.Unmarshal([]byte(jsonify), &dt)

	return GetTimeStampFromTime(dt.RfcDateTime), err
}

func GetTimeStampFromTime(time time.Time) int {
	return int(time.UnixNano() / 1000000)
}

func Float64ToString(num float64) string {
	return strconv.FormatFloat(num, 'f', -1, 64)
}
func Float32ToString(num float64) string {
	return strconv.FormatFloat(num, 'f', -1, 32)
}
