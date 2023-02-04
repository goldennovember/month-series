import pandas as pd
import argparse
import ctypes
from pyarrow.cffi import ffi
import pyarrow.csv as ac
import pyarrow as pa

def read_data(filepath):
    df = pd.read_table(filepath, delimiter=",")
    print(f"Successfully read file from {filepath} with {len(df)} records.")
    return df

def create_table(df, table_name):
    create_statement = pd.io.sql.get_schema(df, table_name)
    return create_statement.encode("utf-8")

def main(data_path: str, table_name: str, library_path: str, function_name: str,num_partitions: str):
    df = read_data(args.data_path)

    # Convert from pandas to arrow
    tbl = pa.Table.from_pandas(df)
    batches = tbl.to_batches(None)

    # Convert from arrow to cffi
    c_schema = ffi.new('struct ArrowSchema*')
    c_array = ffi.new('struct ArrowArray*')
    ptr_schema = int(ffi.cast('uintptr_t', c_schema))
    ptr_array = int(ffi.cast('uintptr_t', c_array))

    # Pass the schema and array pointer
    batches[0].schema._export_to_c(ptr_schema)
    batches[0]._export_to_c(ptr_array)

    create_statement = create_table(df, args.table_name)
    lib = ctypes.CDLL(args.library_path)
    arrowToDB = lib.arrowToDB
    arrowToDB.argtypes = [ctypes.c_longlong, ctypes.c_longlong,ctypes.c_char_p,ctypes.c_char_p,ctypes.c_longlong]
    arrowToDB(ptr_schema, ptr_array, create_statement, table_name.encode("utf-8"), int(num_partitions))


if __name__ == "__main__":
    parser = argparse.ArgumentParser(
        description="Provide arguments to execute data writing to database with Go, Arrow, pgx and COPY."
    )
    parser.add_argument(
        "--data-path",
        "-d",
        dest="data_path",
        required=True,
        help="Filepath to the data that should be written to database.",
    )
    parser.add_argument(
        "--table-name",
        "-t",
        dest="table_name",
        required=True,
        help="Name of the database table that gets created.",
    )
    parser.add_argument(
        "--library-path",
        "-l",
        dest="library_path",
        required=True,
        help="Filepath to the shared C library exported from the Go script.",
    )
    parser.add_argument(
        "--function-name",
        "-f",
        dest="function_name",
        required=False,
        default="arrowToDB",
        help="Name of the function used in the Go script that you would like to call.",
    )
    parser.add_argument(
        "--partitions",
        "-p",
        dest="num_partitions",
        required=False,
        default=1,
        help="Number of data partitions for parallel insertion into the database."
    )


    args = parser.parse_args()
    print(f"Read cli arguments {args}")

    main(**vars(args))


