from functools import cache

def create_keypad(keys, hole_row, hole_col=0):
    def press(current, target):
        current_pos = keys.index(current)
        target_pos = keys.index(target)
        row_diff = (target_pos // 3) - (current_pos // 3)
        col_diff = (target_pos % 3) - (current_pos % 3)

        col_move = "<>"[col_diff > 0] * abs(col_diff)
        row_move = "^v"[row_diff > 0] * abs(row_diff)

        if target_pos // 3 == hole_row and current_pos % 3 == hole_col:
            return col_move + row_move
        elif current_pos // 3 == hole_row and target_pos % 3 == hole_col:
            return row_move + col_move
        else:
            if "<" in col_move:
                return col_move + row_move
            else:
                return row_move + col_move
    return press

# 7  8  9     0  1  2
# 4  5  6     3  4  5
# 1  2  3     6  7  8
#    0  A     x  10 11
numeric = "789456123_0A"
press_numeric = create_keypad(numeric, 3)

#   ^ A     x 1 2
# < v >     3 4 5
directional = "_^A<v>"
press_directional = create_keypad(directional, 0)

def press_keypads_recursive(code, key_funcs):
    levels = len(key_funcs) - 1

    @cache
    def num_presses(current, target, level):
        sequence = key_funcs[level](current, target) + "A"
        if level == levels:
            return len(sequence)
        else:
            length = 0
            c = "A"
            for t in sequence:
                length += num_presses(c, t, level + 1)
                c = t
            return length

    length = 0
    current = "A"
    for target in code:
        length += num_presses(current, target, 0)
        current = target
    return length

def day21(filename, part2=False):
    with open(filename, "r") as f:
        codes = f.read().splitlines()

    nchain = 25 if part2 else 2
    keypad_chain = [press_numeric] + [press_directional] * nchain
    complexity = 0
    for code in codes:
        sequence_length = press_keypads_recursive(code, keypad_chain)
        numeric_code = int("".join(c for c in code if c.isnumeric()))
        complexity += sequence_length * numeric_code
        print(numeric_code, sequence_length, complexity)
    return complexity

if __name__ == "__main__":
    # print("Part 1 example", day21("input/day21_example.txt"))
    print("Part 1", day21("input/input.txt"))
    # print("Part 2 example", day21("input/day21_example.txt", True))
    # print("Part 2", day21("input/day21.txt", True))
