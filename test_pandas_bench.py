# - Group by aggregate
# - Apply

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
])
def test_filter(benchmark, name, filter_fn, expected_count):
    df = read_csv()
    new_df = benchmark(filter_fn, df)
    assert len(new_df) == expected_count


# Name (time in ms)                                             Mean              Median            StdDev            Rounds
# --------------------------------------------------------------------------------------------------------------------------
# test_filter[combine and-<lambda>-7280]                      5.0941 (1.0)        4.9476 (1.0)      0.7620 (1.66)        175
# test_filter[combine or-<lambda>-39818]                     10.0231 (1.97)       9.9354 (2.01)     0.4589 (1.0)         102
# test_filter[contains case insensitive-<lambda>-11912]      98.2903 (19.29)     98.7691 (19.96)    2.2099 (4.82)         10
# test_filter[contains case sensitive-<lambda>-9118]         59.1438 (11.61)     58.6890 (11.86)    3.4401 (7.50)         18
# test_filter[regex case insensitive-<lambda>-11912]        305.6966 (60.01)    306.0521 (61.86)    4.8669 (10.61)         5
# test_filter[regex case sensitive-<lambda>-9118]           128.5158 (25.23)    127.9879 (25.87)    3.7973 (8.27)          9
# test_filter[single float-<lambda>-26823]                    7.4126 (1.46)       7.2683 (1.47)     0.7895 (1.72)        135
# test_filter[string eq-<lambda>-830]                        10.0484 (1.97)       9.8872 (2.00)     0.9117 (1.99)         93
# test_read_csv                                             416.8440 (81.83)    416.2372 (84.13)    2.7442 (5.98)          5
# test_sort[columns0]                                        22.8061 (4.48)      22.5379 (4.56)     1.0598 (2.31)         42
# test_sort[columns1]                                       142.2413 (27.92)    143.0750 (28.92)    2.4375 (5.31)          8
# test_sort[columns2]                                       184.1537 (36.15)    183.4041 (37.07)    5.8602 (12.77)         6
# test_write_json_records                                   317.4334 (62.31)    318.0416 (64.28)    4.9470 (10.78)         5
# --------------------------------------------------------------------------------------------------------------------------
