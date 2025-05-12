#include <SDL3/SDL.h>
#include <SDL3/SDL_main.h>
#include "Game.h"
#include <iostream>
using namespace std;

int main(int argc, char *argv[]) {
  Game game;
  bool isReady = game.Initialize();

  if (isReady) {
    game.RunLoop();
  }

  game.Shutdown();
  std::cout << "SDL Initialized and Quit successfully!\n";
  return 0;
}
