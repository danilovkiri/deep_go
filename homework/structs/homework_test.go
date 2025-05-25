package main

import (
	"fmt"
	"math"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

type Option func(*GamePerson)

func WithName(name string) func(*GamePerson) {
	return func(person *GamePerson) {
		copy(person.name[:], name)
	}
}

func WithCoordinates(x, y, z int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.x = int32(x)
		person.y = int32(y)
		person.z = int32(z)
	}
}

func WithGold(gold int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.gold = uint32(gold)
	}
}

func WithMana(mana int) func(*GamePerson) {
	return func(person *GamePerson) {
		mask := uint32(0x3FF << 22)
		person.attributes2 = (person.attributes2 &^ mask) | (uint32(mana) << 22)
	}
}

func WithHealth(health int) func(*GamePerson) {
	return func(person *GamePerson) {
		mask := uint32(0x3FF << 12)
		person.attributes2 = (person.attributes2 &^ mask) | (uint32(health) << 12)
	}
}

func WithRespect(respect int) func(*GamePerson) {
	return func(person *GamePerson) {
		mask := uint16(0xF << 12)
		person.attributes1 = (person.attributes1 &^ mask) | (uint16(respect) << 12)
	}
}

func WithStrength(strength int) func(*GamePerson) {
	return func(person *GamePerson) {
		mask := uint16(0xF << 8)
		person.attributes1 = (person.attributes1 &^ mask) | (uint16(strength) << 8)
	}
}

func WithExperience(experience int) func(*GamePerson) {
	return func(person *GamePerson) {
		mask := uint16(0xF << 4)
		person.attributes1 = (person.attributes1 &^ mask) | (uint16(experience) << 4)
	}
}

func WithLevel(level int) func(*GamePerson) {
	return func(person *GamePerson) {
		mask := uint16(0xF << 0)
		person.attributes1 = (person.attributes1 &^ mask) | (uint16(level) << 0)
	}
}

func WithHouse() func(*GamePerson) {
	return func(person *GamePerson) {
		mask := uint32(0b1 << 11)
		person.attributes2 = (person.attributes2 &^ mask) | (1 << 11)
	}
}

func WithGun() func(*GamePerson) {
	return func(person *GamePerson) {
		mask := uint32(0b1 << 10)
		person.attributes2 = (person.attributes2 &^ mask) | (1 << 10)
	}
}

func WithFamily() func(*GamePerson) {
	return func(person *GamePerson) {
		mask := uint32(0b1 << 9)
		person.attributes2 = (person.attributes2 &^ mask) | (1 << 9)
	}
}

func WithType(personType int) func(*GamePerson) {
	return func(person *GamePerson) {
		mask := uint32(0b11 << 7)
		person.attributes2 = (person.attributes2 &^ mask) | (uint32(personType) << 7)
	}
}

const (
	BuilderGamePersonType = iota
	BlacksmithGamePersonType
	WarriorGamePersonType
)

type GamePerson struct {
	name        [42]byte // 42 bytes total, 42 characters
	attributes1 uint16   // 2 bytes total, 4 bits for each: respect, strength, experience, level (in this order L->R)
	// mana, health - 10 bits each ([0,1000])
	// home, weapon, family - 1 bit each (boolean)
	// player type - 2 bits ([0,2]) - all in this order L->R
	// total 25 bits -> ceil total 32 bits -> ceil total 4 bytes
	attributes2 uint32
	gold        uint32 // 4 bytes total
	x, y, z     int32  // 12 bytes total
}

func NewGamePerson(options ...Option) GamePerson {
	entity := GamePerson{}
	for idx := range options {
		options[idx](&entity)
	}
	return entity
}

func (p *GamePerson) Name() string {
	return unsafe.String(&p.name[0], len(p.name))
}

func (p *GamePerson) X() int {
	return int(p.x)
}

func (p *GamePerson) Y() int {
	return int(p.y)
}

func (p *GamePerson) Z() int {
	return int(p.z)
}

func (p *GamePerson) Gold() int {
	return int(p.gold)
}

func (p *GamePerson) Mana() int {
	return int(p.attributes2 >> 22 & 0x3FF)
}

func (p *GamePerson) Health() int {
	return int(p.attributes2 >> 12 & 0x3FF)
}

func (p *GamePerson) Respect() int {
	return int(p.attributes1 >> 12 & 0xF)
}

func (p *GamePerson) Strength() int {
	return int(p.attributes1 >> 8 & 0xF)
}

func (p *GamePerson) Experience() int {
	return int(p.attributes1 >> 4 & 0xF)
}

func (p *GamePerson) Level() int {
	return int(p.attributes1 >> 0 & 0xF)
}

func (p *GamePerson) HasHouse() bool {
	return (p.attributes2 >> 11 & 0b1) == 1
}

func (p *GamePerson) HasGun() bool {
	return (p.attributes2 >> 10 & 0b1) == 1
}

func (p *GamePerson) HasFamily() bool {
	return (p.attributes2 >> 9 & 0b1) == 1
}

func (p *GamePerson) Type() int {
	return int(p.attributes2 >> 7 & 0b11)
}

func TestGamePerson(t *testing.T) {
	assert.LessOrEqual(t, unsafe.Sizeof(GamePerson{}), uintptr(64))

	const x, y, z = math.MinInt32, math.MaxInt32, 0
	const name = "aaaaaaaaaaaaa_bbbbbbbbbbbbb_cccccccccccccc"
	const personType = BuilderGamePersonType
	const gold = math.MaxInt32
	const mana = 1000
	const health = 1000
	const respect = 10
	const strength = 10
	const experience = 10
	const level = 10

	options := []Option{
		WithName(name),
		WithCoordinates(x, y, z),
		WithGold(gold),
		WithMana(mana),
		WithHealth(health),
		WithRespect(respect),
		WithStrength(strength),
		WithExperience(experience),
		WithLevel(level),
		WithHouse(),
		WithFamily(),
		WithType(personType),
	}

	person := NewGamePerson(options...)
	assert.Equal(t, name, person.Name())
	assert.Equal(t, x, person.X())
	assert.Equal(t, y, person.Y())
	assert.Equal(t, z, person.Z())
	assert.Equal(t, gold, person.Gold())
	assert.Equal(t, mana, person.Mana())
	assert.Equal(t, health, person.Health())
	assert.Equal(t, respect, person.Respect())
	assert.Equal(t, strength, person.Strength())
	assert.Equal(t, experience, person.Experience())
	assert.Equal(t, level, person.Level())
	assert.True(t, person.HasHouse())
	assert.True(t, person.HasFamily())
	assert.False(t, person.HasGun())
	assert.Equal(t, personType, person.Type())

	fmt.Println(person.Type())
}
