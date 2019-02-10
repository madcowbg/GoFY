# GoFY
Experiments in financial instrument pricing in Go.

# API
This developed as a pure Go library, but JSON-based interface is built-in so that other languages can call particular functionalities.

The current coverage contains option pricing, calculation of implied volatilities and curve bootstrap. It is relatively straightforward
to extend.

# Algorithms
There are several algorithms currently implemented.

## Derivatives pricing
Currently only vanila options and their greeks are supported, in a variety of pricing algorithms:
 * Binomial tree
 * Finite difference schema
 * Pricing by Monte Carlo simulations


## Interest rate curve bootstrap.
Two algorithms are provided for interest rate bootstrap, both based on bond yields.

### Naive interpolation
Uses the bond yields to bootstrap the curve incrementally. Values are matched exactly, and tenors provided for whatever
maturities bonds are given. Constant rates are assumed between points.

Naturally this produces an OK spot curve, but a very discontinuous forward curve with clear arbitrage. Also, tenors are not constant.

### Monotone convex interpolation
This interpolation method makes spot interpolation to be aware of arbitrage and is made to produce no-arbitrage interpolation. 
See bibiliography for math details.

This method is used in bootstrap by OLS on bond prices. A set of requested tenor points is provided to the algorithm.

Note: The initial guess is a naive bootstrap via direct interpolation, which is then constantly interpolated to get the initial guess
at the requested tenors.

# Bibliography
## Derivatives pricing
The derivative pricing algorithms mostly follow Paul Wilmott's book:

Paul Wilmott Introduces Quantitative Finance 2 
Wiley-Interscience New York, NY, USA ©2007 
ISBN:0470319585 9780470319581.

## Interest rate curve interpolation
Monotone Convex interpolation: Hagan, Patrick S., and Graeme West. "Methods for constructing a yield curve." Wilmott Magazine, May (2008): 70-81.
