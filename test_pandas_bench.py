# - Read CSV
# - Write JSON
# - Sort
# - Group by aggregate
# - Filter, string, regex, LIKE, int, float
# - Apply

import pandas as pd


def read_csv():
    return pd.read_csv('recipeData-utf8.csv')


def test_read_csv(benchmark):
    # benchmark something
    df = benchmark(read_csv)
    assert len(df) == 73861
    assert len(list(df)) == 23

# ----------------------------------------------- benchmark: 1 tests -----------------------------------------------
# Name (time in ms)          Min       Max      Mean  StdDev    Median     IQR  Outliers     OPS  Rounds  Iterations
# ------------------------------------------------------------------------------------------------------------------
# test_read_csv         412.2500  425.9627  418.2159  4.9514  417.3080  4.6171       2;0  2.3911       5           1
# ------------------------------------------------------------------------------------------------------------------
