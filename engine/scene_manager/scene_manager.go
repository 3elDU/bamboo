/*
Scene is a distinct state of the program, that displays specific state of the game.
For example: main menu scene, "new game" scene, playing scene, death scene, etc.
*/
package scene_manager

import (
	"fmt"
	"log"
	"reflect"

	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/exp/slices"
)

// Keeping one and global instance of scene manager
var manager *sceneManager

type Scene interface {
	Update() error
	Draw(screen *ebiten.Image)

	// called when the scene is about to be deleted
	Destroy()
}

type sceneManager struct {
	currentScene Scene
	queue        []Scene

	// tick counter
	// can be retrieved throuch Ticks() function,
	// and used for timing purposes, etc.
	counter uint64

	// special flag, that is set in SceneManager.Exit()
	terminated bool
}

func init() {
	ebiten.SetWindowClosingHandled(true)
	manager = &sceneManager{
		currentScene: nil,
		queue:        make([]Scene, 0),
	}
}

// Returns internal tick counter, that is incremented on each Update() call
// Can be used for different timing purposes, etc.
func Ticks() uint64 {
	return manager.counter
}

// Must be called from Scene.Update()
// Exits current scene, and switches to next in the queue
// If the queue is empty, exits
func End() {
	if manager.currentScene != nil {
		manager.currentScene.Destroy()
	}

	if len(manager.queue) != 0 {
		next := manager.queue[0]
		manager.currentScene = next

		// delete scene from the queue
		manager.queue[0] = nil
		manager.queue = slices.Delete(manager.queue, 0, 1)
	} else {
		manager.currentScene = nil
	}

	manager.printQueue("End")
}

// Terminates the program, destroying all the remaining scenes
func Exit() {
	log.Println("SceneManager.Exit() called. Terminating all the scenes and quiting")
	for _, scene := range manager.queue {
		scene.Destroy()
	}
	if manager.currentScene != nil {
		manager.currentScene.Destroy()
	}
	manager.terminated = true
}

// Switches to the given scene, inserting current scene to the queue
// Switch is intented for temporary scenes, like pause menu
func Switch(next Scene) {
	if manager.currentScene != nil {
		manager.queue = slices.Insert(manager.queue, 0, manager.currentScene)
	}

	manager.currentScene = next
	manager.printQueue("Switch")
}

// Behaves similarly to Switch, but the main difference is,
// QSwitch completely replaces current scene with new one
func QSwitch(next Scene) {
	manager.currentScene = next
	manager.printQueue("QSwitch")
}

// Pushes scene to the end of the queue
func Push(sc Scene) {
	manager.queue = append(manager.queue, sc)
	manager.printQueue("Push")
}

// Inserts scene at the beginning of the queue
func Insert(sc Scene) {
	manager.queue = slices.Insert(manager.queue, 0, sc)
	manager.printQueue("Insert")
}

func (manager *sceneManager) Update() error {
	if manager.terminated {
		return fmt.Errorf("exit")
	}

	if ebiten.IsWindowBeingClosed() {
		log.Println("SceneManager.Update() - Handling window close")

		if manager.currentScene != nil {
			manager.currentScene.Destroy()
		}
		for _, sc := range manager.queue {
			sc.Destroy()
		}

		return fmt.Errorf("exit")
	}

	if manager.currentScene == nil {
		if len(manager.queue) != 0 {
			End()
		} else {
			log.Println("SceneManager.Update() - No scenes left to display. Exiting!")
			return fmt.Errorf("exit")
		}
	}

	if err := manager.currentScene.Update(); err != nil {
		return err
	}

	manager.counter++

	return nil
}

func (manager *sceneManager) Draw(screen *ebiten.Image) {
	if manager.currentScene != nil {
		manager.currentScene.Draw(screen)
	}
}

func (manager *sceneManager) Layout(outsideWidth, outsideHeight int) (int, int) {
	return outsideWidth, outsideHeight
}

func (manager *sceneManager) printQueue(originFunc string) {
	queueTypes := make([]reflect.Type, len(manager.queue))
	for i, scene := range manager.queue {
		queueTypes[i] = reflect.TypeOf(scene)
	}

	log.Printf("SceneManager.%v - current %v; queue %v",
		originFunc, reflect.TypeOf(manager.currentScene), queueTypes)
}

func Run() {
	if err := ebiten.RunGame(manager); err != nil {
		switch err.Error() {
		case "exit":
			break
		default:
			log.Panicln(err)
		}
	}
}
