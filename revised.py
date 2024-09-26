import os
import sys
import pygame
from pygame.locals import *
from pynput import keyboard

pygame.init()
WIDTH, HEIGHT = 560, 450
screen = pygame.display.set_mode((WIDTH, HEIGHT))
pygame.display.set_caption("Fractal Flux Origin")

WHITE = (255, 255, 255)
BLACK = (0, 0, 0)
RED = (255, 0, 0)
GRAY = (200, 200, 200)
LIGHT_GRAY = (220, 220, 220)

WHITE_KEY_WIDTH = WIDTH // 7
WHITE_KEY_HEIGHT = 300
BLACK_KEY_WIDTH = WHITE_KEY_WIDTH // 2
BLACK_KEY_HEIGHT = WHITE_KEY_HEIGHT * 2 // 3

POT_CENTER = (WIDTH // 2, 400)
POT_RADIUS = 30
POT_ANGLE_RANGE = 270
pot_val = 1

SLIDER_WIDTH = WIDTH - 40
SLIDER_HEIGHT = 30
SLIDER_X = 20
SLIDER_Y = 320
selector_val = 1

white_keys = ['z', 'x', 'c', 'v', 'b', 'n', 'm']
black_keys = ['s', 'd', 'g', 'h', 'j']

pipe = None
pressed_keys = set()
pot_value = 50  
selector_value = 1  

def setup_pipe():
    global pipe
    pipe_path = "/tmp/pipe_frequency"
    if not os.path.exists(pipe_path):
        os.mkfifo(pipe_path)
    pipe = open(pipe_path, 'w')

def cleanup_pipe():
    global pipe
    if pipe:
        pipe.close()

def write_to_pipe(message):
    global pipe
    if pipe:
        pipe.write(message)
        pipe.flush()

def on_press(key):
    global pressed_keys
    try:
        char = key.char
    except AttributeError:
        char = str(key)

    if char not in pressed_keys:
        pressed_keys.add(char)
        write_to_pipe(f"p:{char}\n")
    
    return char != 'q'

def on_release(key):
    global pressed_keys
    try:
        char = key.char
    except AttributeError:
        char = str(key)

    if char in pressed_keys:
        pressed_keys.remove(char)
        write_to_pipe(f"r:{char}\n")
    
    return char != 'q'

def draw_piano():
    
    for i, key in enumerate(white_keys):
        x = i * WHITE_KEY_WIDTH
        color = RED if key in pressed_keys else WHITE
        pygame.draw.rect(screen, color, (x, 0, WHITE_KEY_WIDTH, WHITE_KEY_HEIGHT))
        pygame.draw.rect(screen, BLACK, (x, 0, WHITE_KEY_WIDTH, WHITE_KEY_HEIGHT), 2)
    
    black_key_positions = [0, 1, 3, 4, 5]
    for i, key in enumerate(black_keys):
        x = (black_key_positions[i] * WHITE_KEY_WIDTH) + (WHITE_KEY_WIDTH * 3 // 4)
        color = RED if key in pressed_keys else BLACK
        pygame.draw.rect(screen, color, (x, 0, BLACK_KEY_WIDTH, BLACK_KEY_HEIGHT))

def draw_potentiometer():
    global pot_value
    pygame.draw.circle(screen, GRAY, POT_CENTER, POT_RADIUS)
    angle = (pot_value / 100) * POT_ANGLE_RANGE - 225
    end_x = POT_CENTER[0] + POT_RADIUS * pygame.math.Vector2(1, 0).rotate(angle).x
    end_y = POT_CENTER[1] + POT_RADIUS * pygame.math.Vector2(1, 0).rotate(angle).y
    pygame.draw.line(screen, BLACK, POT_CENTER, (end_x, end_y), 2)

def draw_selector():
    global selector_value
    pygame.draw.rect(screen, LIGHT_GRAY, (SLIDER_X, SLIDER_Y, SLIDER_WIDTH, SLIDER_HEIGHT))
    pygame.draw.rect(screen, BLACK, (SLIDER_X, SLIDER_Y, SLIDER_WIDTH, SLIDER_HEIGHT), 2)
    
    handle_x = SLIDER_X + (selector_value - 1) * (SLIDER_WIDTH - 20) / 31
    pygame.draw.rect(screen, BLACK, (handle_x, SLIDER_Y, 20, SLIDER_HEIGHT))

def update_potentiometer(y_diff):
    global pot_value, pot_val
    pot_value -= y_diff // 2  
    pot_value = max(0, min(100, pot_value))
    if pot_value != pot_val:
        write_to_pipe(f"pot:{pot_value}\n")
        pot_val = pot_value

def update_selector(x):
    global selector_value, selector_val
    new_value = max(1, min(32, int((x - SLIDER_X) / SLIDER_WIDTH * 32) + 1))
    if new_value != selector_value:
        selector_value = new_value
        if selector_value != selector_val:
            write_to_pipe(f"sel:{selector_value}\n")
            selector_val = selector_value

if __name__ == "__main__":
    try:
        setup_pipe()
        print("Keyboard listener starting...")
        listener = keyboard.Listener(on_press=on_press, on_release=on_release)
        listener.start()
        
        running = True
        dragging_pot = False
        dragging_selector = False
        last_y = 0
        while running:
            screen.fill(WHITE)
            for event in pygame.event.get():
                if event.type == pygame.QUIT:
                    running = False
                elif event.type == pygame.MOUSEBUTTONDOWN:
                    if pygame.math.Vector2(event.pos[0] - POT_CENTER[0], event.pos[1] - POT_CENTER[1]).length() <= POT_RADIUS:
                        dragging_pot = True
                        last_y = event.pos[1]
                    elif SLIDER_X <= event.pos[0] <= SLIDER_X + SLIDER_WIDTH and SLIDER_Y <= event.pos[1] <= SLIDER_Y + SLIDER_HEIGHT:
                        dragging_selector = True
                        update_selector(event.pos[0])
                elif event.type == pygame.MOUSEBUTTONUP:
                    dragging_pot = False
                    dragging_selector = False
                elif event.type == pygame.MOUSEMOTION:
                    if dragging_pot:
                        y_diff = event.pos[1] - last_y
                        update_potentiometer(y_diff)
                        last_y = event.pos[1]
                    elif dragging_selector:
                        update_selector(event.pos[0])

            draw_piano()
            draw_selector()
            draw_potentiometer()
            pygame.display.flip()

    except Exception as e:
        print(f"An error occurred: {e}", file=sys.stderr)
    finally:
        cleanup_pipe()
        pygame.quit()
        print("Keyboard listener stopped.")