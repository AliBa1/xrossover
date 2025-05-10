#include <SDL3/SDL.h>
#include <SDL3/SDL_main.h>
#include <iostream>
using namespace std;

const int screenWidth = 800;
const int screenHeight = 600;

class Game {
public:
  // Game();

  bool Initialize();
  void RunLoop();
  void Shutdown();

private:
  SDL_Window *window;
  SDL_Renderer *renderer;
  bool isRunning;

  void ProcessInput();
  void UpdateGame();
  void GenerateOutput();
};

bool Game::Initialize() {
  if (SDL_Init(SDL_INIT_VIDEO) != true) {
    SDL_Log("SDL_INIT Error: %s", SDL_GetError());
    return false;
  }

  window = SDL_CreateWindow("Game Window", screenWidth, screenHeight, 0);
  if (!window) {
    SDL_Log("SDL_CreateWindow Error: %s", SDL_GetError());
    return false;
  }

  renderer = SDL_CreateRenderer(window, NULL);
  if (!renderer) {
    SDL_Log("SDL_CreateRenderer Error: %s", SDL_GetError());
    return false;
  }
  
  isRunning = true;
  return true;
}

void Game::RunLoop() {
  while (isRunning) {
    ProcessInput();
    GenerateOutput();
    SDL_Delay(100);
  }
}

void Game::Shutdown() {
  SDL_DestroyWindow(window);
  SDL_DestroyRenderer(renderer);
  SDL_Quit();
}

void Game::ProcessInput() {
  SDL_Event event;
  while (SDL_PollEvent(&event)) {
    switch (event.type) {
    case SDL_EVENT_QUIT:
      isRunning = false;
      break;
    }
  }
}

void Game::GenerateOutput() {
  SDL_SetRenderDrawColor(renderer, 0, 0, 255, 255);
  SDL_RenderClear(renderer);
  SDL_RenderPresent(renderer);
}

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
