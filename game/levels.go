package game

import (
	"time"

	"github.com/gdamore/tcell"
	"github.com/google/logger"
	"github.com/stjiub/gosnake/entity"
	"github.com/stjiub/gosnake/gamemap"
)

// Generate level 1 map which is just an open map with walls around perimeter
func InitLevel1(g *Game) {
	i := entity.NewItem(MapWidth/2+3, MapHeight/2+3, WallPass, (time.Second * 3), '*', g.DefStyle)
	g.items = append(g.items, i)
	go randomBits(g, 2, 10, 3*time.Second)
	go randomLines(g, 2)
}

func InitLevel2(g *Game) {
	g.gameMap.BitChan = make(chan bool, 2)
	go g.moveBits(m)
}

func InitLevel3(g *Game) {
	bChan := make(chan bool, 2)
	g.gameMap.BiteChan = append(g.gameMap.BiteChan, bChan)
	go randomBites(g, 1, 3, (20 * time.Second), false, g.gameMap.BiteChan[0])
}

func InitLevel4(g *Game) {
	makeWallChan(g, 2)
	go movingWall(g, 1+15, g.gameMap.Height/4, entity.DirLeft, 2, 15, WallRune, g.DefStyle, g.gameMap.WallChan[0])
	go movingWall(g, g.gameMap.Width-15, (g.gameMap.Height - g.gameMap.Height/4), entity.DirRight, 2, 15, WallRune, g.DefStyle, g.gameMap.WallChan[1])
}

func InitLevel5(g *Game) {
	bChan := make(chan bool, 2)
	g.gameMap.BiteChan = append(g.gameMap.BiteChan, bChan)
	go randomBites(g, 1, 3, (20 * time.Second), true, g.gameMap.BiteChan[1])
}

func InitLevel6(g *Game) {
	g.gameMap.BiteChan[0] <- true
	makeWallChan(g, 4)
	go movingWall(g, g.gameMap.Width/4, 6, entity.DirUp, 1, 7, WallRune, g.DefStyle, g.gameMap.WallChan[2])
	go movingWall(g, (g.gameMap.Width/4 + 1), 6, entity.DirUp, 1, 7, WallRune, g.DefStyle, g.gameMap.WallChan[3])
	go movingWall(g, (g.gameMap.Width - g.gameMap.Width/4), g.gameMap.Height-6, entity.DirDown, 1, 7, WallRune, g.DefStyle, g.gameMap.WallChan[4])
	go movingWall(g, ((g.gameMap.Width - g.gameMap.Width/4) - 1), g.gameMap.Height-6, entity.DirDown, 1, 7, WallRune, g.DefStyle, g.gameMap.WallChan[5])
}

func InitLevel7(g *Game) {
	for i := range g.gameMap.WallChan {
		g.gameMap.WallChan[i] <- true
	}
	g.gameMap.BitChan <- true
}

func randomLines(g *Game, numTimes int) {
	defer func() {
		if r := recover(); r != nil {
			logger.Errorf("Error in RandomLines goroutine: %v", r)
		}
	}()
	//for i := 0; i < numTimes; i++ {
	for {
		g.bits = entity.NewRandomBitLine(g.bits, m, 10, BitRune, g.BitStyle)
		time.Sleep(15 * time.Second)
	}
}

func randomBits(g *Game, bitsGen, bitsMax int, dur time.Duration) {
	defer func() {
		if r := recover(); r != nil {
			logger.Errorf("Error in RandomBits goroutine: %v", r)
		}
	}()
	for {
		for i := 0; i < bitsGen; i++ {
			if len(g.bits)-bitsGen < bitsMax {
				newB := entity.NewRandomBit(m, 10, BitRune, g.BitStyle)
				g.bits = append(g.bits, newB)
			}
		}
		time.Sleep(dur)
	}
}

func randomBites(g *Game, bitesGen, bitesMax int, dur time.Duration, random bool, biteChan chan bool) {
	defer func() {
		if r := recover(); r != nil {
			logger.Errorf("Error in RandomBites goroutine: %v", r)
		}
	}()
	for {
		select {
		default:
			for i := 0; i < bitesGen; i++ {
				if len(g.bites)-bitesGen < bitesMax {
					newB := entity.NewRandomBite(m, BiteRunes, g.BiteExplodedStyle, random)
					g.bites = append(g.bites, newB)
				}
			}
			time.Sleep(dur)
		case <-biteChan:
			return
		}
	}
}

func movingWall(g *Game, x, y, direction, speed, segments int, char rune, style tcell.Style, quit chan bool) {
	e := entity.NewEntity(x, y, direction, speed, char, style)
	e.AddSegment(segments, char, style)
	g.entities = append(g.entities, e)

	defer func() {
		if r := recover(); r != nil {
			logger.Errorf("Error in MovingWall goroutine: %v", r)
		}
	}()
	for {
		select {
		default:
			dx, dy := e.CheckDirection()
			if e.IsBlockedByMap(m, dx, dy) {
				var newPos []*gamemap.Object
				for i := 0; i < e.GetLength(); i++ {
					o := e.GetSegment(e.GetLength() - 1 - i)
					newPos = append(newPos, o)
				}
				e.NewPos(newPos)
				dir := e.GetDirection()
				switch dir {
				case entity.DirUp:
					e.SetDirection(entity.DirDown)
				case entity.DirDown:
					e.SetDirection(entity.DirUp)
				case entity.DirLeft:
					e.SetDirection(entity.DirRight)
				case entity.DirRight:
					e.SetDirection(entity.DirLeft)
				}
			} else {
				e.Move(dx, dy)
			}
			time.Sleep(g.moveInterval(e.GetSpeed(), e.GetDirection()))
		case <-quit:
			return
		}
	}
}

func makeWallChan(g *Game, num int) {
	for i := 0; i < num; i++ {
		c := make(chan bool, 2)
		g.gameMap.WallChan = append(g.gameMap.WallChan, c)
	}
}
