# WH3 DraftBot #

This is a draft bot for use in tournaments in the multiplayer game Warhammer 3.

In drafts, players take turns picking factions for some number of rounds according to a predefined ruleset.

The project is an opportunity to learn more Go, which as of writing I'm not familiar with - so be nice and feedback
welcome ;)

# How does it work #

The bot uses minimax much as Chess engines do to find optimal picks assuming that your opponent also makes optimal picks.

# How to Use #

TODO - Need to sort out CLI

# Formats #

## 2022-Q2-Turin-Default ##
The bot currently supports only the most popular format where if you represent player 1 and player 2 as P1/P2, a bo3 goes:

G1: P1: Pick 2, P2: Pick 1, P1: Pick 1 from the original 2.

G2: P2: Pick 2, P1: Pick 1, P2 pick 1 from the original 2.

G3: Winner G2 Pick 3 ban 1, Other player: Pick 1 ban 1, Winner G2 pick 1 ban 1.

No repeat final picks are allowed but repeat initial picks are permitted.

Bo5/7 are played in the same way.
