# _bamboo_

### What is it?
It is a little game, that I'm working on, as a hobby. I started it when I was learning Go.  
It is a top-down sandbox game, where you are stranded on an island, and need to survive. Craft, kill, build, make yourself home.  
Most of ideas and game mechanics were inspired by [Minecraft](https://minecraft.net/).  

### How to play it?
Clone the repository, and run `go run .` in the terminal.   
[Go](https://go.dev/) must be installed for this to work, of course.

### Progress
- [x] World generation
    - [x] Trees
    - [x] Mushrooms
    - [ ] Biomes
    - [ ] Caves
- [x] The map actually looks like a giant island
- [x] Player movement
    - WASD
- [ ] Player physics
- [x] Map scaling
    - Move your mouse wheel to zoom
- [x] UI
- [x] World saving/loading
- [x] Placing blocks
    - Hold F to place stone blocks under you
- [ ] Inventory
- [ ] Items
- [ ] Crafting
- [ ] Animals
- [ ] Farming

### Plans for the future
- [ ] Multiplayer
- [ ] Lua scripting

## Naming conventions
You may notice, that some functions are suffixed with `B` (example: `world.ChunkAtB`), and some are not.  
The reason is pretty simple: functions with suffix `B` accept block coordinates, and functions without that suffix accept chunk coordinates.