# - Read CSV
# - Write JSON
# - Sort
# - Group by aggregate
# - Filter, string, regex, LIKE, int, float
# - Apply

import pandas as pd


def read_csv():
    return pd.read_csv('recipeData-utf8.csv')


def write_json(df):
    return df.to_json(orient='records')


def test_read_csv(benchmark):
    df = benchmark(read_csv)
    assert len(df) == 73861
    assert len(list(df)) == 23


def test_write_json_records(benchmark):
    df = read_csv()
    data = benchmark(write_json, df)
    assert len(data) == 33565370

# ------------------------------------------------------------------------------------- benchmark: 2 tests -------------------------------------------------------------------------------------
# Name (time in ms)                Min                 Max                Mean            StdDev              Median               IQR            Outliers     OPS            Rounds  Iterations
# ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------
# test_write_json_records     307.5048 (1.0)      322.5679 (1.0)      316.4568 (1.0)      6.2581 (1.04)     319.3597 (1.0)      9.6468 (1.09)          1;0  3.1600 (1.0)           5           1
# test_read_csv               407.9583 (1.33)     423.1593 (1.31)     415.5637 (1.31)     5.9940 (1.0)      417.4608 (1.31)     8.8643 (1.0)           2;0  2.4064 (0.76)          5           1
# ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------
