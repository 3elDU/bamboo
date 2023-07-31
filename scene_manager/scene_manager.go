package scene_manager

import (
	"fmt"
	"log"
	"reflect"

	"github.com/3elDU/bamboo/config"
	"github.com/3elDU/bamboo/event"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"golang.org/x/exp/slices"
)

// Keeping one and global instance of scene manager
var manager *sceneManager

/*
Scene is a distinct state of the program, that displays specific state of the game.
For example: main menu scene, "new game" scene, playing scene, death scene, etc.
*/
type Scene interface {
	Update()
	Draw(screen *ebiten.Image)

	// Destroy is called when the scene is about to be deleted
	Destroy()
}

type sceneManager struct {
	currentScene Scene
	stack        []Scene

	// tick counter
	// can be retrieved through Ticks() function
	counter uint64

	// special flag, that is set in SceneManager.Exit()
	terminated bool
}

func init() {
	ebiten.SetWindowClosingHandled(true)
	manager = &sceneManager{
		currentScene: nil,
		stack:        make([]Scene, 0),
	}
}

// Ticks Returns internal tick counter that is incremented on each Update() call
// Can be used for different timing purposes
func Ticks() uint64 {
	return manager.counter
}

// Pop must be called from Scene.Update()
// Exits current scene, and switches to next in the stack
// If the stack is empty, exits
func Pop() {
	if manager.currentScene != nil {
		manager.currentScene.Destroy()
	}

	if len(manager.stack) != 0 {
		next := manager.stack[0]
		manager.currentScene = next

		// delete scene from the stack
		manager.stack[0] = nil
		manager.stack = slices.Delete(manager.stack, 0, 1)
	} else {
		manager.currentScene = nil
	}

	manager.printQueue("Pop")
}

// Exit terminates the program, destroying all the remaining scenes
func Exit() {
	log.Println("SceneManager.Exit() called. Terminating all the scenes and quiting")
	for _, scene := range manager.stack {
		scene.Destroy()
	}
	if manager.currentScene != nil {
		manager.currentScene.Destroy()
	}
	manager.terminated = true
}

// PushAndSwitch switches to the given scene, inserting current scene to the stack.
// PushAndSwitch is intented for temporary scenes, like pause menu
func PushAndSwitch(next Scene) {
	if manager.currentScene != nil {
		manager.stack = slices.Insert(manager.stack, 0, manager.currentScene)
	}

	manager.currentScene = next
	manager.printQueue("PushAndSwitch")
}

// ReplaceAndSwitch Behaves similarly to PushAndSwitch, but the main difference is,
// ReplaceAndSwitch completely replaces current scene with new one
func ReplaceAndSwitch(next Scene) {
	manager.currentScene = next
	manager.printQueue("QPushAndSwitch")
}

// Push pushes scene to the stack
func Push(sc Scene) {
	manager.stack = append(manager.stack, sc)
	manager.printQueue("Push")
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
		for _, sc := range manager.stack {
			sc.Destroy()
		}

		return fmt.Errorf("exit")
	}

	if manager.currentScene == nil {
		if len(manager.stack) != 0 {
			Pop()
		} else {
			log.Println("SceneManager.Update() - No scenes left to display. Exiting!")
			return fmt.Errorf("exit")
		}
	}

	// Ctrl+F5 to reload everything
	if ebiten.IsKeyPressed(ebiten.KeyControlLeft) && inpututil.IsKeyJustPressed(ebiten.KeyF5) {
		log.Println("Reloading!")
		event.FireEvent(event.NewEvent(event.Reload, nil))
	}

	// F3 to toggle debug mode on/off
	if inpututil.IsKeyJustPressed(ebiten.KeyF3) {
		config.DebugMode = !config.DebugMode
	}

	manager.currentScene.Update()

	manager.counter++

	return nil
}

func (manager *sceneManager) Draw(screen *ebiten.Image) {
	if manager.currentScene != nil {
		manager.currentScene.Draw(screen)
	}
}

func (manager *sceneManager) Layout(outsideWidth, outsideHeight int) (int, int) {
	scaleFactor := ebiten.DeviceScaleFactor()
	return int(scaleFactor * float64(outsideWidth)), int(scaleFactor * float64(outsideHeight))
}

func (manager *sceneManager) printQueue(originFunc string) {
	queueTypes := make([]reflect.Type, len(manager.stack))
	for i, scene := range manager.stack {
		queueTypes[i] = reflect.TypeOf(scene)
	}

	log.Printf("SceneManager.%v - current %v; stack %v",
		originFunc, reflect.TypeOf(manager.currentScene), queueTypes)
}

func Run() {
	if err := ebiten.RunGame(manager); err != nil {
		switch err.Error() {
		// When the window was closed by the user
		case "exit":
			break
		default:
			log.Panicln(err)
		}
	}
}
