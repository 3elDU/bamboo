# _bamboo_

### What is it?
It is a little game, that I'm working on, as a hobby. I started it when I was learning Go.  
It is a top-down sandbox game, where you are stranded on an island, and need to survive. Craft, kill, build, make yourself home.  
Most of ideas and game mechanics were inspired by [Minecraft](https://minecraft.net/).  
The game is in early development, so do not expect stability and backwards-compatibility for world saves.

### How to play it?
Clone the repository, and run `go run .` in the terminal.   
[Go](https://go.dev/) must be installed for this to work, of course.

### Progress
[here](https://github.com/users/3elDU/projects/1)

## Naming conventions
You may notice, that some functions are suffixed with `B` (example: `world.ChunkAtB`), and some are not.  
The reason is pretty simple: functions with suffix `B` accept block coordinates, and functions without that suffix accept chunk coordinates.
