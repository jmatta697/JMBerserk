Author: Joe Matta

JMBerserk - WIZARD BERSERK

(This game is an analogue of the classic Atari 2600 game, Berzerk.)

You are a white mage apprentice. Your attack spell is relatively weak and you can not move very fast. However,
since you're a white mage, you have two uses of a life spell that will revive you if you are killed.

INTRODUCTION:

The object of the game is to destroy all of the dark mages using your only attack spell.
There are four levels, each one has more dark mages and less cover than the last . Also,
the dark mages get more aggressive as you advance into the next levels. You can not advance to the next
level until all dark mages in the that level are destroyed. If you destroy all mages in level four and make it
out alive, you win! Be careful though, the walls are magically charged and will kill anything that touches them.

GAME PLAY:

At the beginning of the game, the player (white mage) is placed on the board at a random location. There are also
4 + (2 * level number) dark mages (the enemy) placed randomly on the board at the beginning of each level. The dark
mages will move and shoot their attack spells DIRECTLY at the player. The frequency of firing and the specific mage
that fires is determined in a random fashion. The player may move in eight directions: up, down, left, right, and all
diagonal directions using the keyboard arrow keys. The player will shoot a projectile spell in the direction of their
last move when the space bar is pressed. There are two ways to destroy the dark mages. One way is hitting them with the
attack spell projectile, and the other way is louring them into a wall. Twenty points are awarded for each dark mage
that is destroyed whether they are killed by the player or by the walls. The current score is displayed din the upper
right hand corner of the game window. The current level is displayed in the upper left hand corner. If any of the dark
mages touch the wall they are destroyed; this also applies to the player. When all of the dark mages for the current
level are destroyed, the player must exit the level to move on to the next level by going through the opening at the
top of the window. When all the mages are destroyed in level four and the player exits level four through the opening,
the game is won; a end game screen will appear. The player can lose a life in three ways: touching any wall,
touching a dark mage, or getting hit by a dark mage projectile. The player has three lives. If all three lives have
expired, the game ends and a 'game over' screen appears.

RUNNING THE GAME:

** Go version 1.11 is required to compile and run this game **

From the terminal -
navigate to the directory ./JMBerserk
when inside the directory run the following command:

$ go run JMBerserk.go

-----
packages needed:

go get..

github.com/faiface/beep
github.com/faiface/pixel
github.com/faiface/mainthread
github.com/faiface/glhf
github.com/hajimehoshi/oto
github.com/go-gl/mathgl/mgl32

golang.org/x/image

you must also run the following commands to get required dev headers:

$ apt install libasound2-dev

(The previous installation may require root (sudo) access to install)

The bellow install may also be required:

$ pkg install openal-soft
