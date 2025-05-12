#include "Game.h"

const int screenWidth = 800;
const int screenHeight = 600;
const float thickness = 15.0;
const float paddleHeight = 15.0 * 7;

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

  ballPos.x = screenWidth / 2.0;
  ballPos.y = screenHeight / 2.0;
  ballVelocity.x = -200.0f;
  ballVelocity.y = 235.0f;

  paddlePos.x = 100;
  paddlePos.y = screenHeight / 2.0;

  ticksCount = 0;
  isRunning = true;
  return true;
}

void Game::RunLoop() {
  while (isRunning) {
    ProcessInput();
    UpdateGame();
    GenerateOutput();
    SDL_Delay(20);
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

  const bool *state = SDL_GetKeyboardState(NULL);

  paddleDir = 0;
  if (state[SDL_SCANCODE_W]) {
    paddleDir -= 1;
  }

  if (state[SDL_SCANCODE_S]) {
    paddleDir += 1;
  }
}

void Game::UpdateGame() {
  // wait for 16ms to limit to 60 FPS
  while ((SDL_GetTicksNS() / 1000000ULL) < (ticksCount + 16));
  // difference in ticks from last frame
  float deltaTime = (SDL_GetTicks() - ticksCount) / 1000.0f;
  ticksCount = SDL_GetTicks();

  if (paddleDir != 0) {
    paddlePos.y += paddleDir * 300.0f * deltaTime;
    // keep paddle on screen
    if (paddlePos.y > (screenHeight - paddleHeight / 2.0f - thickness)) {
      paddlePos.y = screenHeight - paddleHeight / 2.0f - thickness;
    } else if (paddlePos.y < paddleHeight / 2.0f + thickness) {
      paddlePos.y = paddleHeight / 2.0f + thickness;
    }
  }

  ballPos.x += ballVelocity.x * deltaTime;
  ballPos.y += ballVelocity.y * deltaTime;

  float diff = paddlePos.y - ballPos.y;
  diff = (diff > 0.0f) ? diff : -diff;
  // ball hit paddle
  if (
    diff <= paddleHeight / 2.0f &&
    ballPos.x <= paddlePos.x + thickness &&
    ballPos.x >= paddlePos.x &&
    ballVelocity.x < 0.0f
  ) {
    ballVelocity.x *= -1.0f;
  }

  // ball goes off screen
  if (ballPos.x <= 0.0f) {
    isRunning = false;
  }

  // ball collide with right wall
  if (ballPos.x >= (screenWidth - thickness) && ballVelocity.x > 0.0f) {
    ballVelocity.x *= -1;
  }
  // ball collide with top wall
  if (ballPos.y <= thickness && ballVelocity.y < 0.0f) {
    ballVelocity.y *= -1;
  }

  // ball collide with bottom wall
  if (ballPos.y >= (screenHeight - thickness) && ballVelocity.y > 0.0f) {
    ballVelocity.y *= -1;
  }
}

void Game::GenerateOutput() {
  SDL_SetRenderDrawColor(renderer, 0, 0, 255, 255);
  SDL_RenderClear(renderer);

  SDL_SetRenderDrawColor(renderer, 255, 255, 255, 255);

  // top wall
  SDL_FRect wall{0.0, 0.0, static_cast<float>(screenWidth), thickness};
  SDL_RenderFillRect(renderer, &wall);

  // bottom wall
  wall.y = screenHeight - thickness;
  SDL_RenderFillRect(renderer, &wall);

  // right wall
  wall.x = static_cast<float>(screenWidth) - thickness;
  wall.y = 0;
  wall.w = thickness;
  wall.h = static_cast<float>(screenWidth);
  SDL_RenderFillRect(renderer, &wall);

  SDL_FRect ball{ballPos.x - thickness / 2, ballPos.y - thickness / 2,
                 thickness, thickness};
  SDL_RenderFillRect(renderer, &ball);

  // SDL_FRect paddle{paddlePos.x - thickness / 2,
  //                  paddlePos.y - (thickness / 2) - (paddleHeight / 2),
  //                  thickness, paddleHeight};

  SDL_FRect paddle{paddlePos.x,
                   paddlePos.y - (paddleHeight / 2),
                   thickness, paddleHeight};
 
  SDL_RenderFillRect(renderer, &paddle);
  SDL_RenderPresent(renderer);
}
