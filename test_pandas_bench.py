# - Group by aggregate
# - Apply
# - Filter with custom function

import pandas as pd
import pytest


def read_csv():
    return pd.read_csv('recipeData-utf8.csv')


def write_json(df):
    return df.to_json(orient='records')


def sort_df(df, columns):
    return df.sort_values(by=columns)


def test_read_csv(benchmark):
    df = benchmark(read_csv)
    assert len(df) == 73861
    assert len(list(df)) == 23


def test_write_json_records(benchmark):
    df = read_csv()
    data = benchmark(write_json, df)
    assert len(data) == 33565370


@pytest.mark.parametrize("columns", [
    (["UserId"]),
    (["Name"]),
    (["Style", "Name", "BrewMethod"]),
])
def test_sort(benchmark, columns):
    df = read_csv()
    new_df = benchmark(sort_df, df, columns)
    assert len(df) == len(new_df)


@pytest.mark.parametrize("name, filter_fn, expected_count", [
    ("single float", lambda df: df[df["Size(L)"] > 21.0], 26823),
    ("combine or", lambda df: df[(df["Size(L)"] > 21.0) | (df["StyleID"] > 100)], 39818),
    ("combine and", lambda df: df[(df["Size(L)"] > 21.0) & (df["StyleID"] > 100)], 7280),
    ("string eq", lambda df: df[df.Style == "Cream Ale"], 830),
    ("contains case sensitive", lambda df: df[df["Name"].str.contains("Ale", case=True, na=False)], 9118),
    ("contains case insensitive", lambda df: df[df["Name"].str.contains("ale", case=False, na=False)], 11912),
    ("regex case sensitive", lambda df: df[df["Name"].str.contains(".*Ale.*", case=True, na=False)], 9118),
    ("regex case insensitive", lambda df: df[df["Name"].str.contains(".*ale.*", case=False, na=False)], 11912),
    ("integer in", lambda df: df[df["StyleID"].isin([i for i in range(100)])], 53514),
])
def test_filter(benchmark, name, filter_fn, expected_count):
    df = read_csv()
    new_df = benchmark(filter_fn, df)
    assert len(new_df) == expected_count

# Name (time in ms)                                             Mean              Median            StdDev            Rounds
# --------------------------------------------------------------------------------------------------------------------------
# test_filter[combine and-<lambda>-7280]                      4.6996 (1.0)        4.5539 (1.0)      0.4474 (1.0)         192
# test_filter[combine or-<lambda>-39818]                     10.2313 (2.18)      10.0739 (2.21)     0.6902 (1.54)         89
# test_filter[contains case insensitive-<lambda>-11912]     100.4382 (21.37)     99.3378 (21.81)    3.4898 (7.80)         11
# test_filter[contains case sensitive-<lambda>-9118]         59.0800 (12.57)     58.2704 (12.80)    2.0572 (4.60)         18
# test_filter[integer in-<lambda>-53514]                     11.6831 (2.49)      11.4531 (2.51)     0.9151 (2.05)         85
# test_filter[regex case insensitive-<lambda>-11912]        304.9066 (64.88)    306.2898 (67.26)    2.9705 (6.64)          5
# test_filter[regex case sensitive-<lambda>-9118]           126.5499 (26.93)    125.9932 (27.67)    3.2897 (7.35)          8
# test_filter[single float-<lambda>-26823]                    7.3997 (1.57)       7.1886 (1.58)     0.8638 (1.93)        113
# test_filter[string eq-<lambda>-830]                        10.1308 (2.16)       9.9458 (2.18)     1.0313 (2.31)         92
# test_read_csv                                             416.8440 (81.83)    416.2372 (84.13)    2.7442 (5.98)          5
# test_sort[columns0]                                        22.8061 (4.48)      22.5379 (4.56)     1.0598 (2.31)         42
# test_sort[columns1]                                       142.2413 (27.92)    143.0750 (28.92)    2.4375 (5.31)          8
# test_sort[columns2]                                       184.1537 (36.15)    183.4041 (37.07)    5.8602 (12.77)         6
# test_write_json_records                                   317.4334 (62.31)    318.0416 (64.28)    4.9470 (10.78)         5
# --------------------------------------------------------------------------------------------------------------------------
