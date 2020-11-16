package main

import (
	"fmt"
	_ "fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	_ "image/png"
	"log"
	"math/rand"
	"time"
)

const (
	ScreenWidth  = 800
	ScreenHeight = 700
)

type Sprite struct {
	upPict    *ebiten.Image
	downPict  *ebiten.Image
	leftPict  *ebiten.Image
	rightPict *ebiten.Image
	xLoc      int
	yLoc      int
	dx        int
	dy        int
	width     float64
	height    float64
	collision bool
	direction string
	health    int
}

type Game struct {
	playerSprite               Sprite
	personEnemy                Sprite
	monsterEnemy               Sprite
	tankTopper                 Sprite
	fireball                   Sprite
	coinSprite                 Sprite
	heartSprite1               Sprite
	heartSprite2               Sprite
	heartSprite3               Sprite
	firstMap                   Sprite
	secondMap                  Sprite
	thirdMap                   Sprite
	drawOps                    ebiten.DrawImageOptions
	collectedGold              bool
	playerAndWallCollision     bool
	projectileAndWallCollision bool
	mostRecentKeyLeft          bool
	mostRecentKeyRight         bool
	mostRecentKeyDown          bool
	mostRecentKeyUp            bool
	mostRecentKeyA             bool
	mostRecentKeyS             bool
	mostRecentKeyD             bool
	mostRecentKeyW             bool
	deathCounter               int
	score                      int
	projectileList             []Sprite
	levelOneEnemyList          []Sprite
	levelTwoEnemyList          []Sprite
	levelThreeEnemyList        []Sprite
	projectileHold             bool
	levelOneIsActive           bool
	levelTwoIsActive           bool
	levelThreeIsActive         bool
	spawnedLevel1Enemies       bool
	spawnedLevel2Enemies       bool
	spawnedLevel3Enemies       bool
	enemyCanShoot              bool //to prevent enemies from shooting the instant the game starts
	gameOver                   bool
}

func gotGold(player, gold Sprite) bool {
	goldWidth, goldHeight := gold.upPict.Size()
	playerWidth, playerHeight := player.upPict.Size()
	if player.xLoc < gold.xLoc+goldWidth &&
		player.xLoc+playerWidth > gold.xLoc &&
		player.yLoc < gold.yLoc+goldHeight &&
		player.yLoc+playerHeight > gold.yLoc {
		return true
	}
	return false
}

func wallCollisionCheckFirstLevel(anySprite Sprite, spriteWidth int) bool {
	boundaryWidth := 25
	if anySprite.xLoc < 0+boundaryWidth || anySprite.xLoc > ScreenWidth-boundaryWidth-spriteWidth ||
		anySprite.yLoc > ScreenHeight-boundaryWidth-spriteWidth || anySprite.yLoc < 0+boundaryWidth ||
		anySprite.xLoc > 200-spriteWidth && anySprite.xLoc < 275 && anySprite.yLoc < 250 ||
		anySprite.xLoc > 275-spriteWidth && anySprite.xLoc < 475 && anySprite.yLoc < 250 && anySprite.yLoc > 175-spriteWidth ||
		anySprite.xLoc > 175-spriteWidth && anySprite.xLoc < 275 && anySprite.yLoc < 475 && anySprite.yLoc > 400-spriteWidth ||
		anySprite.xLoc > 550-spriteWidth && anySprite.xLoc < 625 && anySprite.yLoc < 575 && anySprite.yLoc > 350-spriteWidth ||
		anySprite.xLoc > 475-spriteWidth && anySprite.xLoc < 550 && anySprite.yLoc < 575 && anySprite.yLoc > 500-spriteWidth {
		return true
	}
	return false
}

func projectileCollisionWithEnemy(anyEnemy Sprite, anyProjectileSprite Sprite, enemyWidth int, projectileWidth int) (bool, bool, int, int) {
	if (anyProjectileSprite.xLoc < anyEnemy.xLoc+enemyWidth &&
		anyProjectileSprite.xLoc+projectileWidth > anyEnemy.xLoc &&
		anyProjectileSprite.yLoc < anyEnemy.yLoc+enemyWidth &&
		anyProjectileSprite.yLoc+projectileWidth > anyEnemy.yLoc) && (anyEnemy.health == 1) {
		anyEnemy.health -= 1
		additionalScore := 200
		return true, true, anyEnemy.health, additionalScore
	} else if (anyProjectileSprite.xLoc < anyEnemy.xLoc+enemyWidth &&
		anyProjectileSprite.xLoc+projectileWidth > anyEnemy.xLoc &&
		anyProjectileSprite.yLoc < anyEnemy.yLoc+enemyWidth &&
		anyProjectileSprite.yLoc+projectileWidth > anyEnemy.yLoc) && (anyEnemy.health == 2) {
		fmt.Println("here")
		anyEnemy.health -= 1
		additionalScore := 100
		return false, true, anyEnemy.health, additionalScore
	}
	additionalScore := 0
	return false, false, anyEnemy.health, additionalScore
}

func (game *Game) shootFireball() []Sprite {
	if inpututil.IsKeyJustReleased(ebiten.KeySpace) && game.projectileHold == false {
		game.projectileHold = true
		go func() {
			<-time.After(500 * time.Millisecond)
			game.projectileHold = false
		}()
		game.projectileAndWallCollision = false
		tempFireball := game.fireball

		if game.mostRecentKeyW == true {
			tempFireball.xLoc = game.playerSprite.xLoc + 20
			tempFireball.yLoc = game.playerSprite.yLoc - 18
			tempFireball.dx = 0
			tempFireball.dy = -10
			game.projectileList = append(game.projectileList, tempFireball)
		} else if game.mostRecentKeyS == true {
			tempFireball.xLoc = game.playerSprite.xLoc + 20
			tempFireball.yLoc = game.playerSprite.yLoc + 55
			tempFireball.dx = 0
			tempFireball.dy = 10
			game.projectileList = append(game.projectileList, tempFireball)
		} else if game.mostRecentKeyA == true {
			tempFireball.xLoc = game.playerSprite.xLoc - 15
			tempFireball.yLoc = game.playerSprite.yLoc + 18
			tempFireball.dx = -10
			tempFireball.dy = 0
			game.projectileList = append(game.projectileList, tempFireball)
		} else if game.mostRecentKeyD == true {
			tempFireball.xLoc = game.playerSprite.xLoc + 55
			tempFireball.yLoc = game.playerSprite.yLoc + 18
			tempFireball.dx = 10
			tempFireball.dy = 0
			game.projectileList = append(game.projectileList, tempFireball)
		} else {
			tempFireball.xLoc = game.playerSprite.xLoc + 20
			tempFireball.yLoc = game.playerSprite.yLoc - 18
			tempFireball.dx = 0
			tempFireball.dy = -10
			game.projectileList = append(game.projectileList, tempFireball)
		}
	}
	return game.projectileList
}

func (game *Game) changeTankDirection() {
	if inpututil.IsKeyJustPressed(ebiten.KeyLeft) {
		game.playerSprite.dx = -3
		game.mostRecentKeyLeft = true
		game.mostRecentKeyDown = false
		game.mostRecentKeyRight = false
		game.mostRecentKeyUp = false
	} else if inpututil.IsKeyJustPressed(ebiten.KeyRight) {
		game.playerSprite.dx = 3
		game.mostRecentKeyLeft = false
		game.mostRecentKeyDown = false
		game.mostRecentKeyRight = true
		game.mostRecentKeyUp = false
	} else if inpututil.IsKeyJustReleased(ebiten.KeyRight) || inpututil.IsKeyJustReleased(ebiten.KeyLeft) {
		game.playerSprite.dx = 0
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyUp) {
		game.playerSprite.dy = -3
		game.mostRecentKeyLeft = false
		game.mostRecentKeyDown = false
		game.mostRecentKeyRight = false
		game.mostRecentKeyUp = true
	} else if inpututil.IsKeyJustPressed(ebiten.KeyDown) {
		game.playerSprite.dy = 3
		game.mostRecentKeyLeft = false
		game.mostRecentKeyDown = true
		game.mostRecentKeyRight = false
		game.mostRecentKeyUp = false
	} else if inpututil.IsKeyJustReleased(ebiten.KeyUp) || inpututil.IsKeyJustReleased(ebiten.KeyDown) {
		game.playerSprite.dy = 0
	}
	game.playerSprite.yLoc += game.playerSprite.dy
	game.playerSprite.xLoc += game.playerSprite.dx
}

func (game *Game) changeTankTopperDirection() {
	if inpututil.IsKeyJustPressed(ebiten.KeyW) {
		game.mostRecentKeyA = false
		game.mostRecentKeyS = false
		game.mostRecentKeyD = false
		game.mostRecentKeyW = true
	} else if inpututil.IsKeyJustPressed(ebiten.KeyS) {
		game.mostRecentKeyA = false
		game.mostRecentKeyS = true
		game.mostRecentKeyD = false
		game.mostRecentKeyW = false
	} else if inpututil.IsKeyJustPressed(ebiten.KeyA) {
		game.mostRecentKeyA = true
		game.mostRecentKeyS = false
		game.mostRecentKeyD = false
		game.mostRecentKeyW = false
	} else if inpututil.IsKeyJustPressed(ebiten.KeyD) {
		game.mostRecentKeyA = false
		game.mostRecentKeyS = false
		game.mostRecentKeyD = true
		game.mostRecentKeyW = false
	}
}
func (game *Game) manageTankTopperOffset() {
	if game.mostRecentKeyA == true {
		game.tankTopper.xLoc = game.playerSprite.xLoc - 7
		game.tankTopper.yLoc = game.playerSprite.yLoc + 20
	} else if game.mostRecentKeyD == true {
		game.tankTopper.xLoc = game.playerSprite.xLoc + 20
		game.tankTopper.yLoc = game.playerSprite.yLoc + 20
	} else if game.mostRecentKeyS == true {
		game.tankTopper.xLoc = game.playerSprite.xLoc + 20
		game.tankTopper.yLoc = game.playerSprite.yLoc + 20
	} else if game.mostRecentKeyW == true {
		game.tankTopper.xLoc = game.playerSprite.xLoc + 20
		game.tankTopper.yLoc = game.playerSprite.yLoc - 7
	} else {
		game.tankTopper.xLoc = game.playerSprite.xLoc + 20
		game.tankTopper.yLoc = game.playerSprite.yLoc - 7
	}
}

func (game *Game) spawnLevel1Enemies() {
	if game.spawnedLevel1Enemies == false {
		personEnemy1 := game.personEnemy
		personEnemy2 := game.personEnemy
		monsterEnemy1 := game.monsterEnemy
		monsterEnemy2 := game.monsterEnemy
		personEnemy1.direction = "down"
		personEnemy2.direction = "left"
		monsterEnemy1.direction = "right"
		monsterEnemy2.direction = "left"
		personEnemy1.health = 1
		personEnemy2.health = 1
		monsterEnemy1.health = 2
		monsterEnemy2.health = 2
		personEnemy1.xLoc = 90
		personEnemy1.yLoc = 40
		personEnemy2.xLoc = 425
		personEnemy2.yLoc = 285
		monsterEnemy1.xLoc = 300
		monsterEnemy1.yLoc = 85
		monsterEnemy2.xLoc = 650
		monsterEnemy2.yLoc = 600

		game.levelOneEnemyList = append(game.levelOneEnemyList, personEnemy1)
		game.levelOneEnemyList = append(game.levelOneEnemyList, personEnemy2)
		game.levelOneEnemyList = append(game.levelOneEnemyList, monsterEnemy1)
		game.levelOneEnemyList = append(game.levelOneEnemyList, monsterEnemy2)
	}
	game.spawnedLevel1Enemies = true
}

func (game *Game) movementLevel1Enemies() {
	personEnemyMovementSpeed := 1
	if len(game.levelOneEnemyList) == 4 {
		for i := 0; i < len(game.levelOneEnemyList); i++ {
			//personEnemy1 moves up and down along left side
			if i == 0 {
				if game.levelOneEnemyList[i].direction == "down" && game.levelOneEnemyList[i].yLoc < 500 {
					game.levelOneEnemyList[i].dy = personEnemyMovementSpeed
					game.levelOneEnemyList[i].yLoc += game.levelOneEnemyList[i].dy
				} else if game.levelOneEnemyList[i].direction == "down" && game.levelOneEnemyList[i].yLoc >= 500 {
					game.levelOneEnemyList[i].direction = "up"
					game.levelOneEnemyList[i].dy = 0
				} else if game.levelOneEnemyList[i].direction == "up" &&
					game.levelOneEnemyList[i].yLoc > 40 {
					game.levelOneEnemyList[i].dy = -personEnemyMovementSpeed
					game.levelOneEnemyList[i].yLoc += game.levelOneEnemyList[i].dy
				} else if game.levelOneEnemyList[i].direction == "up" &&
					game.levelOneEnemyList[i].yLoc <= 40 {
					game.levelOneEnemyList[i].direction = "down"
					game.levelOneEnemyList[i].dy = 0
				}
			} else if i == 1 {
				// personEnemy2 moves around in a square
				if game.levelOneEnemyList[i].direction == "left" && game.levelOneEnemyList[i].yLoc <= 285 &&
					game.levelOneEnemyList[i].xLoc <= 425 && game.levelOneEnemyList[i].xLoc > 285 {
					game.levelOneEnemyList[i].dx = -personEnemyMovementSpeed
					game.levelOneEnemyList[i].xLoc += game.levelOneEnemyList[i].dx
				} else if game.levelOneEnemyList[i].direction == "left" && game.levelOneEnemyList[i].yLoc <= 285 &&
					game.levelOneEnemyList[i].xLoc <= 285 {
					game.levelOneEnemyList[i].direction = "down"
					game.levelOneEnemyList[i].dx = 0
				} else if game.levelOneEnemyList[i].direction == "down" &&
					game.levelOneEnemyList[i].yLoc < 425 {
					game.levelOneEnemyList[i].dy = personEnemyMovementSpeed
					game.levelOneEnemyList[i].yLoc += game.levelOneEnemyList[i].dy
				} else if game.levelOneEnemyList[i].direction == "down" && game.levelOneEnemyList[i].yLoc >= 425 {
					game.levelOneEnemyList[i].direction = "right"
					game.levelOneEnemyList[i].dy = 0
				} else if game.levelOneEnemyList[i].direction == "right" && game.levelOneEnemyList[i].yLoc >= 425 &&
					game.levelOneEnemyList[i].xLoc < 425 {
					game.levelOneEnemyList[i].dx = personEnemyMovementSpeed
					game.levelOneEnemyList[i].xLoc += game.levelOneEnemyList[i].dx
				} else if game.levelOneEnemyList[i].direction == "right" && game.levelOneEnemyList[i].yLoc >= 425 &&
					game.levelOneEnemyList[i].xLoc >= 425 {
					game.levelOneEnemyList[i].direction = "up"
					game.levelOneEnemyList[i].dx = 0
				} else if game.levelOneEnemyList[i].direction == "up" && game.levelOneEnemyList[i].xLoc >= 425 &&
					game.levelOneEnemyList[i].yLoc > 285 {
					game.levelOneEnemyList[i].dy = -personEnemyMovementSpeed
					game.levelOneEnemyList[i].yLoc += game.levelOneEnemyList[i].dy
				} else if game.levelOneEnemyList[i].direction == "up" && game.levelOneEnemyList[i].yLoc <= 285 {
					game.levelOneEnemyList[i].direction = "left"
					game.levelOneEnemyList[i].dy = 0
				}
			} else if i == 2 {
				//monsterEnemy1 moves back and forth left and right at the top and chases if in certain proximity
				if game.levelOneEnemyList[i].direction == "right" && game.levelOneEnemyList[i].xLoc < 600 {
					game.levelOneEnemyList[i].dx = personEnemyMovementSpeed
					game.levelOneEnemyList[i].xLoc += game.levelOneEnemyList[i].dx
				} else if game.levelOneEnemyList[i].direction == "right" && game.levelOneEnemyList[i].xLoc >= 600 {
					game.levelOneEnemyList[i].direction = "left"
					game.levelOneEnemyList[i].dx = 0
				} else if game.levelOneEnemyList[i].direction == "left" &&
					game.levelOneEnemyList[i].xLoc > 300 {
					game.levelOneEnemyList[i].dx = -personEnemyMovementSpeed
					game.levelOneEnemyList[i].xLoc += game.levelOneEnemyList[i].dx
				} else if game.levelOneEnemyList[i].direction == "left" &&
					game.levelOneEnemyList[i].xLoc <= 300 {
					game.levelOneEnemyList[i].direction = "right"
					game.levelOneEnemyList[i].dx = 0
				}
			} else if i == 3 {
				//monsterEnemy2 moves back and forth left and right at the bottom and chases if in certain proximity
				if game.levelOneEnemyList[i].direction == "left" && game.levelOneEnemyList[i].xLoc > 100 {
					game.levelOneEnemyList[i].dx = -personEnemyMovementSpeed
					game.levelOneEnemyList[i].xLoc += game.levelOneEnemyList[i].dx
				} else if game.levelOneEnemyList[i].direction == "left" && game.levelOneEnemyList[i].xLoc <= 100 {
					game.levelOneEnemyList[i].direction = "right"
					game.levelOneEnemyList[i].dx = 0
				} else if game.levelOneEnemyList[i].direction == "right" &&
					game.levelOneEnemyList[i].xLoc < 700 {
					game.levelOneEnemyList[i].dx = personEnemyMovementSpeed
					game.levelOneEnemyList[i].xLoc += game.levelOneEnemyList[i].dx
				} else if game.levelOneEnemyList[i].direction == "right" &&
					game.levelOneEnemyList[i].xLoc >= 700 {
					game.levelOneEnemyList[i].direction = "left"
					game.levelOneEnemyList[i].dx = 0
				}
			}
		}
	}

}

func (game *Game) manageLevel1CollisionDetection() {
	if game.collectedGold == false {
		game.collectedGold = gotGold(game.playerSprite, game.coinSprite)
	}

	if game.playerAndWallCollision == false {
		game.playerAndWallCollision = wallCollisionCheckFirstLevel(game.playerSprite, 61)
	} else {
		game.playerSprite.yLoc = ScreenHeight / 2
		game.playerSprite.xLoc = 74 //player width
		game.playerAndWallCollision = false
		game.deathCounter += 1
	}

	if len(game.levelOneEnemyList) > 0 {
		for i := 0; i < len(game.levelOneEnemyList); i++ {
			if game.levelOneEnemyList[i].collision == false {
				if game.levelOneEnemyList[i].direction == "left" {
					spriteWidth, _ := game.levelOneEnemyList[i].leftPict.Size()
					game.levelOneEnemyList[i].collision = wallCollisionCheckFirstLevel(game.levelOneEnemyList[i], spriteWidth)
				} else if game.levelOneEnemyList[i].direction == "right" {
					spriteWidth, _ := game.levelOneEnemyList[i].rightPict.Size()
					game.levelOneEnemyList[i].collision = wallCollisionCheckFirstLevel(game.levelOneEnemyList[i], spriteWidth)
				} else if game.levelOneEnemyList[i].direction == "up" {
					spriteWidth, _ := game.levelOneEnemyList[i].upPict.Size()
					game.levelOneEnemyList[i].collision = wallCollisionCheckFirstLevel(game.levelOneEnemyList[i], spriteWidth)
				} else if game.levelOneEnemyList[i].direction == "down" {
					spriteWidth, _ := game.levelOneEnemyList[i].downPict.Size()
					game.levelOneEnemyList[i].collision = wallCollisionCheckFirstLevel(game.levelOneEnemyList[i], spriteWidth)
				}
			} else {
				game.levelOneEnemyList[i].dx = 0
				game.levelOneEnemyList[i].dy = 0
			}
		}
	}

	if len(game.projectileList) > 0 {
		for i := 0; i < len(game.projectileList); i++ {
			if game.projectileList[i].collision == false {
				game.projectileList[i].xLoc += game.projectileList[i].dx
				game.projectileList[i].yLoc += game.projectileList[i].dy
				game.projectileList[i].collision = wallCollisionCheckFirstLevel(game.projectileList[i], 20)
			}
		}
	}
	if len(game.projectileList) > 0 && len(game.levelOneEnemyList) > 0 {
		for i := 0; i < len(game.projectileList); i++ {
			for j := 0; j < len(game.levelOneEnemyList); j++ {
				enemyWidth, _ := game.levelOneEnemyList[j].upPict.Size()
				if game.levelOneEnemyList[j].collision == false && game.projectileList[i].collision == false {
					additionalScore := 0
					game.levelOneEnemyList[j].collision, game.projectileList[i].collision, game.levelOneEnemyList[j].health, additionalScore =
						projectileCollisionWithEnemy(game.levelOneEnemyList[j], game.projectileList[i], enemyWidth, 20)
					game.score += additionalScore
				}
			}
		}
	}
}

func (game *Game) checkLevel() {
	if game.gameOver == false {
		if game.score < 1000 {
			game.levelOneIsActive = true
			game.levelTwoIsActive = false
			game.levelThreeIsActive = false
		} else if game.score >= 1000 && game.score < 2000 {
			game.levelOneIsActive = false
			game.levelTwoIsActive = true
			game.levelThreeIsActive = false
		} else if game.score >= 2000 {
			game.levelOneIsActive = false
			game.levelTwoIsActive = false
			game.levelThreeIsActive = true
		} else {
			game.levelOneIsActive = true
			game.levelTwoIsActive = false
			game.levelThreeIsActive = false
		}
	} else {
		game.levelOneIsActive = true
		game.levelTwoIsActive = false
		game.levelThreeIsActive = false
	}
}

func (game *Game) Update() error {
	game.checkLevel()

	if game.deathCounter >= 3 {
		game.gameOver = true
	} else {
		game.gameOver = false
	}
	if len(game.levelOneEnemyList) == 4 {
		for i := 0; i < len(game.levelOneEnemyList); i++ {
		}
	}
	if game.levelOneIsActive == true && game.gameOver == false {
		game.spawnLevel1Enemies()
		game.movementLevel1Enemies()
		game.changeTankDirection()
		game.changeTankTopperDirection()
		game.shootFireball()
		game.manageTankTopperOffset()
		game.manageLevel1CollisionDetection()
		fmt.Println(game.score)

	} else if game.levelTwoIsActive == true && game.gameOver == false {

	} else if game.levelThreeIsActive == true && game.gameOver == false {

	} else if game.gameOver == true {

	} else {
		game.spawnLevel1Enemies()
		game.changeTankDirection()
		game.changeTankTopperDirection()
		game.shootFireball()
		game.manageTankTopperOffset()
		game.manageLevel1CollisionDetection()
	}
	return nil
}

func (game Game) Draw(screen *ebiten.Image) {
	//screen.Fill(colornames.Chocolate)
	if game.gameOver == false {
		game.drawOps.GeoM.Reset()
		game.drawOps.GeoM.Translate(float64(game.firstMap.xLoc), float64(game.firstMap.yLoc))
		screen.DrawImage(game.firstMap.upPict, &game.drawOps)

		if len(game.levelOneEnemyList) > 0 {
			for i := 0; i < len(game.levelOneEnemyList); i++ {
				if game.levelOneEnemyList[i].collision == false {
					game.drawOps.GeoM.Reset()
					game.drawOps.GeoM.Translate(float64(game.levelOneEnemyList[i].xLoc), float64(game.levelOneEnemyList[i].yLoc))
					if game.levelOneEnemyList[i].direction == "left" {
						screen.DrawImage(game.levelOneEnemyList[i].leftPict, &game.drawOps)
					} else if game.levelOneEnemyList[i].direction == "right" {
						screen.DrawImage(game.levelOneEnemyList[i].rightPict, &game.drawOps)
					} else if game.levelOneEnemyList[i].direction == "up" {
						screen.DrawImage(game.levelOneEnemyList[i].upPict, &game.drawOps)
					} else if game.levelOneEnemyList[i].direction == "down" {
						screen.DrawImage(game.levelOneEnemyList[i].downPict, &game.drawOps)
					} else {
						screen.DrawImage(game.levelOneEnemyList[i].upPict, &game.drawOps)
					}
				}
			}
		}

		if len(game.projectileList) > 0 {
			for i := 0; i < len(game.projectileList); i++ {
				if game.projectileList[i].collision == false {
					game.drawOps.GeoM.Reset()
					game.drawOps.GeoM.Translate(float64(game.projectileList[i].xLoc), float64(game.projectileList[i].yLoc))
					screen.DrawImage(game.projectileList[i].upPict, &game.drawOps)
				}
			}
		}

		game.drawOps.GeoM.Reset()
		game.drawOps.GeoM.Translate(float64(game.playerSprite.xLoc), float64(game.playerSprite.yLoc))
		if game.mostRecentKeyUp == true {
			screen.DrawImage(game.playerSprite.upPict, &game.drawOps)
		} else if game.mostRecentKeyDown == true {
			screen.DrawImage(game.playerSprite.downPict, &game.drawOps)
		} else if game.mostRecentKeyRight == true {
			screen.DrawImage(game.playerSprite.rightPict, &game.drawOps)
		} else if game.mostRecentKeyLeft == true {
			screen.DrawImage(game.playerSprite.leftPict, &game.drawOps)
		} else {
			screen.DrawImage(game.playerSprite.upPict, &game.drawOps)
		}

		game.drawOps.GeoM.Reset()
		game.drawOps.GeoM.Translate(float64(game.tankTopper.xLoc), float64(game.tankTopper.yLoc))
		if game.mostRecentKeyW == true {
			screen.DrawImage(game.tankTopper.upPict, &game.drawOps)
		} else if game.mostRecentKeyS == true {
			screen.DrawImage(game.tankTopper.downPict, &game.drawOps)
		} else if game.mostRecentKeyD == true {
			screen.DrawImage(game.tankTopper.rightPict, &game.drawOps)
		} else if game.mostRecentKeyA == true {
			screen.DrawImage(game.tankTopper.leftPict, &game.drawOps)
		} else {
			screen.DrawImage(game.tankTopper.upPict, &game.drawOps)
		}

		if game.deathCounter == 0 {
			game.drawOps.GeoM.Reset()
			game.drawOps.GeoM.Translate(float64(game.heartSprite1.xLoc), float64(game.heartSprite1.yLoc))
			screen.DrawImage(game.heartSprite1.upPict, &game.drawOps)
			game.drawOps.GeoM.Reset()
			game.drawOps.GeoM.Translate(float64(game.heartSprite2.xLoc), float64(game.heartSprite2.yLoc))
			screen.DrawImage(game.heartSprite2.upPict, &game.drawOps)
			game.drawOps.GeoM.Reset()
			game.drawOps.GeoM.Translate(float64(game.heartSprite3.xLoc), float64(game.heartSprite3.yLoc))
			screen.DrawImage(game.heartSprite3.upPict, &game.drawOps)
		} else if game.deathCounter == 1 {
			game.drawOps.GeoM.Reset()
			game.drawOps.GeoM.Translate(float64(game.heartSprite1.xLoc), float64(game.heartSprite1.yLoc))
			screen.DrawImage(game.heartSprite1.upPict, &game.drawOps)
			game.drawOps.GeoM.Reset()
			game.drawOps.GeoM.Translate(float64(game.heartSprite2.xLoc), float64(game.heartSprite2.yLoc))
			screen.DrawImage(game.heartSprite2.upPict, &game.drawOps)
		} else if game.deathCounter == 2 {
			game.drawOps.GeoM.Reset()
			game.drawOps.GeoM.Translate(float64(game.heartSprite1.xLoc), float64(game.heartSprite1.yLoc))
			screen.DrawImage(game.heartSprite1.upPict, &game.drawOps)
		} else if game.deathCounter > 2 {
			game.gameOver = true
		}

		if game.collectedGold == false {
			game.drawOps.GeoM.Reset()
			game.drawOps.GeoM.Translate(float64(game.coinSprite.xLoc), float64(game.coinSprite.yLoc))
			screen.DrawImage(game.coinSprite.upPict, &game.drawOps)
		}
	}
}

func (g Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return ScreenWidth, ScreenHeight
}

func main() {
	ebiten.SetWindowSize(ScreenWidth, ScreenHeight)
	ebiten.SetWindowTitle("Berserk/Tank Game by Trevor Wysong")
	gameObject := Game{}
	loadImage(&gameObject)

	gameObject.tankTopper.xLoc = gameObject.playerSprite.xLoc
	gameObject.tankTopper.yLoc = gameObject.playerSprite.yLoc

	playerWidth, _ := gameObject.playerSprite.upPict.Size()
	gameObject.playerSprite.xLoc = playerWidth
	gameObject.playerSprite.yLoc = ScreenHeight / 2

	coinWidth, coinHeight := gameObject.coinSprite.upPict.Size()
	rand.Seed(int64(time.Now().Second()))
	gameObject.coinSprite.xLoc = rand.Intn(ScreenWidth - coinWidth)
	gameObject.coinSprite.yLoc = rand.Intn(ScreenHeight - coinHeight)

	boundaryWidth := 25
	heartWidth, heartHeight := gameObject.heartSprite1.upPict.Size()
	gameObject.heartSprite1.yLoc = ScreenHeight - (boundaryWidth * 2) - (heartHeight / 2)
	gameObject.heartSprite1.xLoc = boundaryWidth + 16
	gameObject.heartSprite2.yLoc = ScreenHeight - (boundaryWidth * 2) - (heartHeight / 2)
	gameObject.heartSprite2.xLoc = (boundaryWidth + 20) + (heartWidth)
	gameObject.heartSprite3.yLoc = ScreenHeight - (boundaryWidth * 2) - (heartHeight / 2)
	gameObject.heartSprite3.xLoc = (boundaryWidth + 24) + (heartWidth * 2)

	if err := ebiten.RunGame(&gameObject); err != nil {
		log.Fatal("Oh no! something terrible happened", err)
	}
}

func loadImage(game *Game) {
	firstMap, _, err := ebitenutil.NewImageFromFile("Level1Correct.png")
	if err != nil {
		log.Fatal("failed to load image", err)
	}
	game.firstMap.upPict = firstMap

	upPlayer, _, err := ebitenutil.NewImageFromFile("tankFilledTopSquare.png")
	if err != nil {
		log.Fatal("failed to load image", err)
	}
	downPlayer, _, err := ebitenutil.NewImageFromFile("tankFilledTopSquareDown.png")
	if err != nil {
		log.Fatal("failed to load image", err)
	}
	leftPlayer, _, err := ebitenutil.NewImageFromFile("tankFilledTopSquareLeft.png")
	if err != nil {
		log.Fatal("failed to load image", err)
	}
	rightPlayer, _, err := ebitenutil.NewImageFromFile("tankFilledTopSquareRight.png")
	if err != nil {
		log.Fatal("failed to load image", err)
	}
	game.playerSprite.upPict = upPlayer
	game.playerSprite.downPict = downPlayer
	game.playerSprite.leftPict = leftPlayer
	game.playerSprite.rightPict = rightPlayer

	tankTopperUp, _, err := ebitenutil.NewImageFromFile("tankTopper.png")
	if err != nil {
		log.Fatal("failed to load image", err)
	}
	tankTopperDown, _, err := ebitenutil.NewImageFromFile("tankTopperDown.png")
	if err != nil {
		log.Fatal("failed to load image", err)
	}
	tankTopperLeft, _, err := ebitenutil.NewImageFromFile("tankTopperLeft.png")
	if err != nil {
		log.Fatal("failed to load image", err)
	}
	tankTopperRight, _, err := ebitenutil.NewImageFromFile("tankTopperRight.png")
	if err != nil {
		log.Fatal("failed to load image", err)
	}
	game.tankTopper.upPict = tankTopperUp
	game.tankTopper.downPict = tankTopperDown
	game.tankTopper.leftPict = tankTopperLeft
	game.tankTopper.rightPict = tankTopperRight

	fireball, _, err := ebitenutil.NewImageFromFile("fireball.png")
	if err != nil {
		log.Fatal("failed to load image", err)
	}
	game.fireball.upPict = fireball

	coins, _, err := ebitenutil.NewImageFromFile("gold-coins-large.png")
	if err != nil {
		log.Fatal("failed to load image", err)
	}
	game.coinSprite.upPict = coins

	personEnemyUp, _, err := ebitenutil.NewImageFromFile("personEnemyUp.png")
	if err != nil {
		log.Fatal("failed to load image", err)
	}
	personEnemyDown, _, err := ebitenutil.NewImageFromFile("personEnemyDown.png")
	if err != nil {
		log.Fatal("failed to load image", err)
	}
	personEnemyLeft, _, err := ebitenutil.NewImageFromFile("personEnemyLeft.png")
	if err != nil {
		log.Fatal("failed to load image", err)
	}
	personEnemyRight, _, err := ebitenutil.NewImageFromFile("personEnemyRight.png")
	if err != nil {
		log.Fatal("failed to load image", err)
	}
	game.personEnemy.upPict = personEnemyUp
	game.personEnemy.downPict = personEnemyDown
	game.personEnemy.leftPict = personEnemyLeft
	game.personEnemy.rightPict = personEnemyRight

	monsterEnemyUp, _, err := ebitenutil.NewImageFromFile("monsterEnemyUp.png")
	if err != nil {
		log.Fatal("failed to load image", err)
	}
	monsterEnemyDown, _, err := ebitenutil.NewImageFromFile("monsterEnemyDown.png")
	if err != nil {
		log.Fatal("failed to load image", err)
	}
	monsterEnemyLeft, _, err := ebitenutil.NewImageFromFile("monsterEnemyLeft.png")
	if err != nil {
		log.Fatal("failed to load image", err)
	}
	monsterEnemyRight, _, err := ebitenutil.NewImageFromFile("monsterEnemyRight.png")
	if err != nil {
		log.Fatal("failed to load image", err)
	}
	game.monsterEnemy.upPict = monsterEnemyUp
	game.monsterEnemy.downPict = monsterEnemyDown
	game.monsterEnemy.leftPict = monsterEnemyLeft
	game.monsterEnemy.rightPict = monsterEnemyRight

	heart, _, err := ebitenutil.NewImageFromFile("heartScaled.png")
	if err != nil {
		log.Fatal("failed to load image", err)
	}
	game.heartSprite1.upPict = heart
	game.heartSprite2.upPict = heart
	game.heartSprite3.upPict = heart
}
