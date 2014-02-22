package main

import (
	"fmt"
	"math/rand"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func main() {
	args := os.Args[1:]

	if len(args) != 1 {
		fmt.Println("ERROR: must have only one argument")
	} else {
		fmt.Println(rollDice(args[0]))
	}
}

func rollDice(rolls string) int {
	dice := parseRoll(rolls)

	turn := NewRoll()
	for _, die := range dice {
		throws, sides, modifier := parseDie(die)
		turn = turn.die(NewDice(sides), throws, modifier)
	}

	return turn.throw()
}

func parseDie(die string) (int, int, int) {
	parts := strings.Split(die, "d")

	// empty string is "1"
	if parts[0] == "" {
		parts[0] = "1"
	}

	// convert string to integers
	times, err := strconv.ParseInt(parts[0], 0, 0)
	sides, err := strconv.ParseInt(parts[1], 0, 0)
	if err != nil {
		panic(err)
	}

	modifier := int64(1)
	if times < 0 {
		modifier = -1
		times = times * modifier
	}

	return int(times), int(sides), int(modifier)
}

func parseRoll(roll string) []string {
	// re := regexp.MustCompile(`[+-]?\d*d\d+|[+-]{1}m\d+`)
	re := regexp.MustCompile(`[+-]?\d*d\d+`)
	return re.FindAllString(roll, -1)
}

// -------

type dice struct {
	sides int
}

var random = rand.New(rand.NewSource(time.Now().UnixNano()))

func NewDice(sides int) *dice {
	return &dice{sides}
}

func (d *dice) Roll() int {
	return random.Intn(d.sides) + 1
}

func (d *dice) Sides() int {
	return d.sides
}

// --------

type roll struct {
	d        *dice
	throws   int
	modifier int
	nextRoll *roll
}

// NewRoll creates and returns and empty roll.
func NewRoll() *roll {
	return &roll{}
}

func newRollWithDice(die *dice, throws int, modifier int) *roll {
	return &roll{die, throws, modifier, nil}
}

func (r *roll) throw() int {
	if r.d == nil {
		return 0
	}

	var results int
	for i := 0; i < r.throws; i++ {
		res := r.d.Roll() * r.modifier
		results += res
	}

	if r.nextRoll == nil {
		return results
	} else {
		return results + r.nextRoll.throw()
	}
}

func (r *roll) die(die *dice, throws int, modifier int) *roll {
	newRoll := newRollWithDice(die, throws, modifier)
	newRoll.nextRoll = r
	return newRoll
}
