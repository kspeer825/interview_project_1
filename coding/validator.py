import sys
import re


def is_valid(credit_card: str):
    """
    A validator method which checks the format of a credit card satisfies a series of conditions.

    :return: boolean value based on whether or not the credit card format is valid
    """

    # It must start with a 4, 5 or 6.
    if credit_card[0] in ["4", "5", "6"]:

         # It must contain exactly 16 digits.
         # It must only consist of digits (0-9) and hyphens "-".
         # It may have digits in groups of 4, separated by one hyphen "-".
         if not (bool(re.fullmatch(r'\d{16}', credit_card)) or
                 bool(re.fullmatch(r'\d{4}-\d{4}-\d{4}-\d{4}', credit_card))):
             return False

         # It must NOT have 4 or more consecutive repeated digits.
         parsed_credit_card = credit_card.replace('-', '')
         for index in range(10):
             pattern = r'.*' + re.escape(str(index)) + r'{4}.*'
             if bool(re.fullmatch(pattern, parsed_credit_card)):
                 return False

         return True

    return False

def sanitize_input(credit_card_input: list):
    """
    Ensure provided input meets expected format.
    The first line of input contains an integer N such that 0 < N < 100
    The next N lines contain credit card numbers.

    :return: a list of provided credit card values
    """

    if not credit_card_input:
        raise RuntimeError("No input provided to STDIN")

    count = credit_card_input[0]
    credit_cards = credit_card_input[1:]

    valid_counts = [str(n) for n in range(1,100)]
    if count not in valid_counts:
        raise RuntimeError("Invalid input. Line 1 must be an integer N such that 0 < N < 100.")

    return credit_cards

def accepting_input():
    """
    Accept input from STDIN. The first line is excepted to be an integer N indicating the number
    of lines remaining. Read N + 1 lines and return a list of strings.

    :return: list of user input where each item is a line from STDIN
    """

    credit_card_input = []
    for line in sys.stdin:
        credit_card_count = len(credit_card_input) - 1

        # stop accepting input if N + 1 lines have been provided
        if credit_card_count > 1 and credit_card_count == int(credit_card_input[0]):
            break

        credit_card_input.append(line.rstrip())

    return credit_card_input

def main():
    """A credit card validator."""

    credit_card_input = accepting_input()
    credit_cards = sanitize_input(credit_card_input)

    for credit_card in credit_cards:
        validity = "Valid" if is_valid(credit_card) else "Invalid"
        print(validity)


if __name__ == '__main__':
    main()
