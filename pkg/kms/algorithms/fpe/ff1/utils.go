package ff1

import (
	"fmt"
	"math/big"
)

func num(X []int, radix int) (big.Int, error) {
	var r = big.NewInt(int64(radix))
	var x = big.NewInt(0)

	for i := 0; i < len(X); i++ {
		if X[i] < 0 || X[i] > radix {
			return *x, fmt.Errorf("X[%d] is not within the range of values defined by the radix (0..%d)", i, radix)
		}
		// let x = x * radix + X[i]
		var z = big.NewInt(int64(X[i]))
		x = x.Mul(x, r)
		x = x.Add(x, z)
	}
	return *x, nil
}

/**
 * NIST SP 800-38G Algorithm 3: STR<sup>m</sup><sub>radix</sub>(x) -
 * Converts an integer to an array of numerals of a given radix.
 * <p>
 * Prerequisites:<br>
 * Base, radix;<br>
 * String length, m.
 * <p>
 * Input:<br>
 * Integer, x, such that 0 &lt;= x &lt; radix<sup>m</sup>.
 * <p>
 * Output:<br>
 * Numeral string, X.
 *
 * @param x
 *            The integer to convert to a string of numerals.
 * @param radix
 *            The template of the numerals such that 0 &lt;= X[i] &lt; radix for
 *            all i.
 * @param m
 *            The length of the string of numerals.
 */
func str(x big.Int, radix int, m int) ([]int, error) {
	var r = big.NewInt(int64(radix))

	var X = make([]int, m)

	var a big.Int

	for i := 1; i <= m; i++ {

		// x
		a.Set(&x)
		// i. X[m+1-i] = x mod radix;
		z := a.Mod(&a, r).Int64()
		X[m-i] = int(z)

		// ii. x = floor(x/radix).
		x.Div(&x, r)
	}
	return X, nil
}
