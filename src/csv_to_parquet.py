import pandas as pd

def write_parquet_file():
    df = pd.read_csv('recipeData-utf8.csv')
    df.to_parquet('data/recipeData.parquet')

write_parquet_file()
