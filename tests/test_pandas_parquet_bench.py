import pytest

import pandas as pd

def read_parquet():
    return pd.read_parquet('data/recipeData.parquet')


def write_json(df):
    return df.to_json(orient='records')


def sort_df(df, columns):
    return df.sort_values(by=columns)


def test_read_csv(benchmark):
    # benchmark(read_parquet()) ==> cannot get this working!
    df = read_parquet()
    assert len(df) == 73861
    assert len(list(df)) == 23


def test_write_json_records(benchmark):
    df = read_parquet()
    data = benchmark(write_json, df)
    assert len(data) == 33565370


@pytest.mark.parametrize("columns", [
    (["UserId"]),
    (["Name"]),
    (["Style", "Name", "BrewMethod"]),
])
def test_sort(benchmark, columns):
    df = read_parquet()
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
    df = read_parquet()
    new_df = benchmark(filter_fn, df)
    assert len(new_df) == expected_count


def eval_fn(df, expr):
    return df.eval(expr)


@pytest.mark.parametrize("name, eval_expr", [
    ("float abs", "destCol = abs(BoilSize)"),
    ("float add", "destCol = OG + FG"),
])
def test_eval(benchmark, name, eval_expr):
    df = read_parquet()
    new_df = benchmark(eval_fn, df, eval_expr)
    assert len(new_df) == len(df)
    assert len(list(df)) + 1 == len(list(new_df))


@pytest.mark.parametrize("name, aggregation_fn, expected_count", [
    ("single col string single float mean", lambda df: df.groupby(["Style"], as_index=False).agg({'OG': ['mean']}), 175),
    ("single col int single float mean", lambda df: df.groupby(["StyleID"], as_index=False).agg({'OG': ['mean']}), 176),
    ("double col string single float mean", lambda df: df.groupby(["Style", "Color"], as_index=False).agg({'OG': ['mean']}), 39065),
    ("single col string double float mean", lambda df: df.groupby(["Style"], as_index=False).agg({'OG': ['mean'], 'FG': ['mean']}), 175),
])
def test_aggregation(benchmark, name, aggregation_fn, expected_count):
    df = read_parquet()
    new_df = benchmark(aggregation_fn, df)
    assert len(new_df) == expected_count

