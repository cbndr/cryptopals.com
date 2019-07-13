This solution took me a while, partially because the [instructions on cryptopals](https://cryptopals.com/sets/1/challenges/6)
one was a were not really straightforward.

This is what I have done to solve the challenge:

To guess the key size, take pairs of blocks from input data and build a list of "key size candidates" with hamming distance; then sort the candidates by averaged hamming distance (function GuessKeySize)

For the best key size candidates, brute force guess the key (function GuessRepeatingXorKey) by using the solution to challenge 4. For this, build blocks of text with every n-th byte of input data, where n is the key size and use the english letter frequency function to identify the likely XOR value.

Observations:

- this method of breaking the key works best if the encrypted data is longer (actually MUCH longer) than the key
- it also works only when unencrypted data is English text