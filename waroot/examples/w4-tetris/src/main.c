#include "wasm4.h"

#include "utils.h"

const uint8_t PIECE_SPRITES[][8] = {
    // L
    {
        0b11111111,
        0b10000001,
        0b10111101,
        0b10100001,
        0b10100001,
        0b10100001,
        0b10000001,
        0b11111111,
    },

    // J
    {
        0b11111111,
        0b10000001,
        0b10000101,
        0b10000101,
        0b10000101,
        0b10111101,
        0b10000001,
        0b11111111,
    },

    // Z
    {
        0b11111111,
        0b10100101,
        0b11001001,
        0b10010011,
        0b10100101,
        0b11001001,
        0b10010011,
        0b11111111,
    },

    // S
    {
        0b11111111,
        0b10100101,
        0b10010011,
        0b11001001,
        0b10100101,
        0b10010011,
        0b11001001,
        0b11111111,
    },

    // O
    {
        0b11111111,
        0b10000001,
        0b10000001,
        0b10000001,
        0b10000001,
        0b10000001,
        0b10000001,
        0b11111111,
    },

    // T
    {
        0b11111111,
        0b10000001,
        0b10111101,
        0b10100101,
        0b10100101,
        0b10111101,
        0b10000001,
        0b11111111,
    },

    // I
    {
        0b11111111,
        0b10000001,
        0b10100101,
        0b10000001,
        0b10100101,
        0b10111101,
        0b10000001,
        0b11111111,
    },
};

const int PIECE_COLORS[] = {
    0x0484d1, // L: blue
    0xfb922b, // J: orange
    0xe53b44, // Z: red
    0x63c64d, // S: green
    0xffe762, // O: yellow
    0xbf66bf, // T: purple
    0x2ce8f4, // I: cyan
};

#define SPIN_CW 1
#define SPIN_CCW -1

#define BOARD_WIDTH 10
#define BOARD_HEIGHT 20
#define BOARD_OFFSET_X 16
char board[BOARD_WIDTH*BOARD_HEIGHT];

typedef enum {
    L, J, Z, S, O, T, I
} PieceType;

char PIECE_COORDS[] = {
    // L
    1, -1,
    -1, 0,
    0, 0,
    1, 0,

    // J
    -1, -1,
    -1, 0,
    0, 0,
    1, 0,

    // Z
    0, -1,
    -1, -1,
    0, 0,
    1, 0,

    // S
    0, -1,
    1, -1,
    0, 0,
    -1, 0,

    // O
    -1, -1,
    -1, 0,
    0, 0,
    0, -1,

    // T
    0, -1,
    -1, 0,
    0, 0,
    1, 0,

    // I
    -2, 0,
    -1, 0,
    0, 0,
    1, 0,
};

const char LEVEL_SPEED[] = {
    53,
    49,
    45,
    41,
    37,
    33,
    28,
    22,
    17,
    11,
    10,
    9,
    8,
    7,
    6,
    6,
    5,
    5,
    4,
    4,
    3,
};

const short CLEAR_SCORES[] = { 40, 100, 300, 1200 };

typedef struct {
    unsigned char x, y;
    PieceType type;
    char coords[8];
} Piece;
Piece piece;
PieceType nextPieceType;

int level;
int score;
int best;
int lines;

char gravityDelay;

char repeatDelay;
char holdingDirection;
char lastGamepadState;

// How many frames remaining in the row clear animation
char clearAnimationDelay;

char gameover;
char gameoverElapsed;

// Bitmask for rows that are painted a color during the row clear animation
unsigned int clearAnimationRowMask;

void drawBlock (PieceType type, int x, int y) {
    blit(PIECE_SPRITES[type], BOARD_OFFSET_X + x, y, 8, 8, BLIT_1BPP);
}

void drawPiece (int x, int y, PieceType type, const char coords[8]) {
    const char* c = coords;
    for (int n = 0; n < 4; ++n) {
        char cx = *(c++);
        char cy = *(c++);
        drawBlock(type, 8*(x+cx), 8*(y+cy));
    }
}

void nextPiece () {
    piece.x = 5;
    piece.y = 1;
    piece.type = nextPieceType;
    nextPieceType = rand() % 7;
    /* tracef("Piece type %d", piece.type); */

    char* c = &PIECE_COORDS[piece.type*8];
    for (int n = 0; n < 8; ++n) {
        piece.coords[n] = *(c++);
    }

    PALETTE[1] = PIECE_COLORS[piece.type];
}

int checkFilledRows () {
    clearAnimationRowMask = 0;

    for (int y = BOARD_HEIGHT-1; y >= 0; --y) {
        int filled = 1;
        for (int x = 0; x < BOARD_WIDTH; ++x) {
            if (!board[y*BOARD_WIDTH + x]) {
                filled = 0;
                break;
            }
        }
        if (filled) {
            clearAnimationRowMask |= 1 << y;
            clearAnimationDelay += 16;
        }
    }

    if (clearAnimationRowMask) {
        return 1;
    }
    return 0;
}

void clearFilledRows () {
    int cleared = 0;

    for (int y = BOARD_HEIGHT-1; y >= 0; --y) {
        int filled = 1;
        for (int x = 0; x < BOARD_WIDTH; ++x) {
            if (!board[y*BOARD_WIDTH + x]) {
                filled = 0;
                break;
            }
        }
        if (filled) {
            ++cleared;

            for (int n = (y+1)*BOARD_WIDTH-1; n >= BOARD_WIDTH; --n) {
                board[n] = board[n-BOARD_WIDTH];
            }
            ++y;
        }
    }

    if (cleared) {
        score += (level+1) * CLEAR_SCORES[cleared-1];
        if (score >= best) {
            best = score;
            diskw(&best, sizeof(best));
        }
        lines += cleared;
        char newLevel = lines/8;
        if (newLevel != level) {
            level = newLevel;
            tone(130 | (940 << 16), 20, 100, TONE_PULSE2 | TONE_MODE2);
        }
    }
}

void copyPieceToBoard () {
    char* c = piece.coords;
    for (int n = 0; n < 4; ++n) {
        char cx = *(c++);
        char cy = *(c++);
        char x = piece.x+cx, y = piece.y+cy;
        board[y*BOARD_WIDTH + x] = piece.type+1;
    }
}

int movePiece (char dx, char dy) {
    char* c = piece.coords;
    for (int n = 0; n < 4; ++n) {
        char cx = *(c++);
        char cy = *(c++);

        char x = piece.x+cx + dx, y = piece.y+cy + dy;
        if (x < 0 || x >= BOARD_WIDTH || y < 0 || y >= BOARD_HEIGHT || board[y*BOARD_WIDTH+x]) {
            return 1;
        }
    }

    piece.x += dx;
    piece.y += dy;
    return 0;
}

int spinPiece (char direction) {
    if (piece.type == O) {
        return 1;
    }

    char coordsCopy[8];
    for (int n = 0; n < 8; ++n) {
        coordsCopy[n] = piece.coords[n];
    }

    for (int n = 0; n < 4; ++n) {
        char cx = piece.coords[2*n];
        char cy = piece.coords[2*n+1];
        piece.coords[2*n] = -direction * cy;
        piece.coords[2*n+1] = direction * cx;
    }

    if (movePiece(0, 0)) {
        for (int n = 0; n < 8; ++n) {
            piece.coords[n] = coordsCopy[n];
        }
        return 1;
    }

    return 0;
}

int stepGravity () {
    if (movePiece(0, 1)) {
        copyPieceToBoard();
        if (!checkFilledRows()) {
            nextPiece();
            if (movePiece(0, 0)) {
                gameover = 1;
                gameoverElapsed = 0;
                tone(1024 | (10 << 16), 0 | (160 << 8), 100, TONE_PULSE1 | TONE_MODE2);

            } else {
                // Play a thunk
                tone(240, 0 | (12 << 8), 100, TONE_NOISE);
            }
        }
        return 1;
    }
    return 0;
}

void start () {
    PALETTE[0] = 0xd8dee9;
    PALETTE[2] = 0x4c566a;
    PALETTE[3] = 0x2e3440;

    diskr(&best, sizeof(best));

    nextPiece();
}

void update () {
    char gamepad = *GAMEPAD1;
    char pressedThisFrame = gamepad & (gamepad ^ lastGamepadState);
    lastGamepadState = gamepad;

    if (clearAnimationDelay) {
        if (clearAnimationDelay % 16 == 0) {
            tone(140 | (540 << 16), 8 | (16 << 8), 100, TONE_NOISE);
        }

        if (--clearAnimationDelay == 0) {
            clearFilledRows();
            clearAnimationRowMask = 0;
            nextPiece();
        }

    } else if (gameover) {
        if (gameoverElapsed < 80) {
            ++gameoverElapsed;
        }

    } else {
        if (pressedThisFrame & (BUTTON_1 | BUTTON_2 | BUTTON_UP)) {
            // Button 2 (B) is usually left of Button 1 (A).
            // Therefore, Button 1 should spins clockwise (to the right) and
            // Button 2 should spin counter-clockwise (to the left).
            // The default spin direction is clockwise.
            if (!spinPiece(pressedThisFrame & BUTTON_2? SPIN_CCW : SPIN_CW)) {
                tone(210 | (250 << 16), 8, 100, TONE_PULSE1 | TONE_MODE3);
            }
        }

        if (pressedThisFrame & BUTTON_DOWN) {
            // Hard drop
            while (!stepGravity()) { }
        }

        if (gamepad & (BUTTON_RIGHT | BUTTON_LEFT)) {
            if (repeatDelay) {
                --repeatDelay;
            } else {
                repeatDelay = holdingDirection ? 6 : 18;
                if (gamepad & BUTTON_LEFT) {
                    movePiece(-1, 0);
                } else {
                    movePiece(1, 0);
                }
            }
            holdingDirection = 1;

        } else {
            holdingDirection = 0;
            repeatDelay = 0;
        }

        if (gravityDelay) {
            --gravityDelay;
        } else {
            gravityDelay = LEVEL_SPEED[min(level, sizeof(LEVEL_SPEED)-1)];
            stepGravity();
        }
    }

    *DRAW_COLORS = 3;
    rect(0, 0, BOARD_OFFSET_X, 160);
    rect(BOARD_OFFSET_X+BOARD_WIDTH*8, 0, 160-(BOARD_OFFSET_X+BOARD_WIDTH*8), 160);

    *DRAW_COLORS = 4;
    rect(BOARD_OFFSET_X, 0, 1, 160);
    rect(BOARD_OFFSET_X+BOARD_WIDTH*8-1, 0, 1, 160);

    int boardY = 0;
    if (gameover) {
        boardY = 2*gameoverElapsed;

        *DRAW_COLORS = 4;
        text("GAME OVER", BOARD_OFFSET_X + 4, 64);
        text(" PRESS X", BOARD_OFFSET_X + 4, 64+8+2);

        // Restart
        if (boardY >= 160 && (pressedThisFrame & BUTTON_1)) {
            gameover = 0;
            level = 0;
            score = 0;
            lines = 0;
            nextPiece();
            for (int ii = 0; ii < BOARD_WIDTH*BOARD_HEIGHT; ++ii) {
                board[ii] = 0;
            }
        }
    }

    for (int y = 0; y < BOARD_HEIGHT; ++y) {
        *DRAW_COLORS = ((clearAnimationRowMask >> y) & 1) && ((clearAnimationDelay>>3) & 1) ? 0x42 : 0x43;
        for (int x = 0; x < BOARD_WIDTH; ++x) {
            char b = board[y*BOARD_WIDTH+x];
            if (b) {
                drawBlock(b-1, 8*x, 8*y + boardY);
            }
        }
    }

    if (!gameover) {
        *DRAW_COLORS = 0x41;
        drawPiece(13, 14, nextPieceType, PIECE_COORDS + 8*nextPieceType);
    }

    *DRAW_COLORS = 4;

    text("LEVEL", 104, 16);
    text("SCORE", 104, 40);
    text("BEST", 104, 64);

    *DRAW_COLORS = 1;
    text("WASM-4", 112, 152);

    char buffer[256];
    text(itoa(buffer, level+1), 104, 24+2);
    text(itoa(buffer, score), 104, 48+2);
    text(itoa(buffer, best), 104, 72+2);

    if (!clearAnimationDelay && !gameover) {
        *DRAW_COLORS = 0x42;
        drawPiece(piece.x, piece.y, piece.type, piece.coords);
    }
}
