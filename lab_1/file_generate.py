import csv
import random
from decimal import Decimal
from multiprocessing import Pool

# Параметры генерации
num_files = 5
num_rows_per_file = 10**8 # 10**8
categories = ['a', 'b', 'c', 'd']
value_range = (1, 1000)

def generate_file(file_num):
    filename = f"file_{file_num}.csv"
    print(f"Создаю файл {filename}...")

    with open(filename, 'w', newline='') as file:
        writer = csv.writer(file)
        writer.writerow(['Category', 'Value'])

        for _ in range(num_rows_per_file):
            category = random.choice(categories)
            value = round(random.uniform(*value_range), 5)
            writer.writerow([category, value])

    print(f"Файл {filename} создан.")

if __name__ == "__main__":
    # Создаем пул процессов
    with Pool(processes=num_files) as pool:
        pool.map(generate_file, range(1, num_files + 1))
