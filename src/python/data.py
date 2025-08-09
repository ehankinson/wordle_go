from dataclasses import dataclass

@dataclass
class GameConfig:
    correct: str
    present: str
    absent: str
    does_not_exist: str
    row_selector: str
    tile_selector: str
    game_grid: str



NYT_CONFIG = GameConfig(
    correct="correct",
    present="present",
    absent="absent",
    does_not_exist="tbd",
    row_selector="[role='group']:nth-child({row_number})",
    tile_selector="[data-testid='tile']",
    game_grid="[data-testid='game-board']"
)



WORDLY_CONFIG = GameConfig(
    correct="letter-correct",
    present="letter-elsewhere",
    absent="letter-absent",
    does_not_exist="letter-tbd",
    row_selector="div.Row:nth-child({row_number})",
    tile_selector="div.Row-letter",
    game_grid="div.game_rows"
)


GAME_CONFIGS = {
    "nty": NYT_CONFIG,
    "wordly": WORDLY_CONFIG
}







