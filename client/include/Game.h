// #ifndef GAME_H
// #define GAME_H

#pragma once

#include <SDL3/SDL.h>

struct Vector2 {
  float x;
  float y;
};

class Game {
public:
  // Game();

  bool Initialize();
  void RunLoop();
  void Shutdown();

private:
  SDL_Window *window;
  SDL_Renderer *renderer;
  Uint64 ticksCount;
  bool isRunning;

  Vector2 ballPos;
  Vector2 ballVelocity;
  Vector2 paddlePos;
  int paddleDir = 0;

  void ProcessInput();
  void UpdateGame();
  void GenerateOutput();
};

// #endif
