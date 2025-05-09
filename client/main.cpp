#include <SDL3/SDL.h>
#include <SDL3/SDL_main.h>
#include <iostream>

int main(int argc, char *argv[]) {
  if (SDL_Init(SDL_INIT_VIDEO) != true) {
    std::cerr << "SDL_INIT Error: " << SDL_GetError() << std::endl;
    return 1;
  }

  SDL_Window *win = SDL_CreateWindow("Hello SDL3", 640, 480, SDL_WINDOW_OPENGL);
  if (!win) {
    std::cerr << "SDL_CreateWindow Error: " << SDL_GetError() << std::endl;
    SDL_Quit();
    return 1;
  }

  SDL_Delay(2000);
  SDL_DestroyWindow(win);
  SDL_Quit();
  std::cout << "SDL Initialized and Quit successfully!\n";
  return 0;
}
