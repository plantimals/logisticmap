#!/usr/bin/env python

import numpy
import scipy.optimize as optimization

def main():
    xs = numpy.array([1,2,3,4,8])
    ys = numpy.array([4.5, 2.75, 2.17, 1.87, 1.435])
    x0 = numpy.array([0.0, 0.0, 0.0])
    sigma = numpy.array([1.0,1.0,1.0,1.0,1.0,1.0])
    def func(x, a, b, c):
        return a + b*x + c*x*x

    print optimization.curve_fit(func, xs, ys, x0) #, sigma)

if __name__ == "__main__":
    main()
