import matplotlib.pyplot as plt
import numpy as np
import keras
import os

# Kind of like how much "context" each pixel gets
matrixSize: int = 7
adjacents: int = int(matrixSize/2)


def load_training_pixel_data():
    pixels_file = open("training_pixel_data.dat", "br")
    stride = int.from_bytes(pixels_file.read(2), "big")
    pixels_raw = pixels_file.read()
    pixels = [bytearray(pixels_raw[i:i+3]) for i in range(0, len(pixels_raw), 3)]
    pixel_matrix = [(pixels[i:i+stride]) for i in range(0, len(pixels), stride)]
    tensors: list[list[list]] = []
    for y in range(len(pixel_matrix)):
        for x in range(len(pixel_matrix[y])):
            tensor: list[list] = []
            for yr in range(y - adjacents, y + adjacents + 1):
                for xr in range(x - adjacents, x + adjacents + 1):
                    if yr < 0 or yr >= len(pixel_matrix) or xr < 0 or xr >= len(pixel_matrix[yr]):
                        tensor.append([255, 255, 255])
                    else:
                        tensor.append(list(pixel_matrix[yr][xr]))
            tensors.append(tensor)
    return np.array(tensors)

def load_prediction_pixel_data():
    pixels_file = open("prediction_pixel_data.dat", "br")
    stride = int.from_bytes(pixels_file.read(2), "big")
    pixels_raw = pixels_file.read()
    pixels = [bytearray(pixels_raw[i:i+3]) for i in range(0, len(pixels_raw), 3)]
    pixel_matrix = [(pixels[i:i+stride]) for i in range(0, len(pixels), stride)]
    tensors: list[list[list]] = []
    for y in range(len(pixel_matrix)):
        for x in range(len(pixel_matrix[y])):
            tensor: list[list] = []
            for yr in range(y - adjacents, y + adjacents + 1):
                for xr in range(x - adjacents, x + adjacents + 1):
                    if yr < 0 or yr >= len(pixel_matrix) or xr < 0 or xr >= len(pixel_matrix[yr]):
                        tensor.append([255, 255, 255])
                    else:
                        tensor.append(list(pixel_matrix[yr][xr]))
            tensors.append(tensor)
    return np.array(tensors)

def load_training_label_data():
    labels_file = open("training_label_data.dat", "br")
    labels_bytes = labels_file.read()
    return np.array(list(labels_bytes))

    
print("Processing pixel training data...")
training_pixel_data = load_training_pixel_data()
print(training_pixel_data.shape)
print("Processing training label data...")
training_label_data = load_training_label_data()
print("Processing prediction data...")
prediction_pixel_data = load_prediction_pixel_data()

training_pixel_data = training_pixel_data / 255.0

model = keras.Sequential([
    keras.layers.Flatten(input_shape=(matrixSize*matrixSize, 3)),
    keras.layers.Dense(128, activation='relu'),
    keras.layers.Dense(5)
])

model.compile(optimizer='adam',
              loss=keras.losses.SparseCategoricalCrossentropy(from_logits=True),
              metrics=['accuracy'])

model.fit(training_pixel_data, training_label_data, epochs=10)

probability_model = keras.Sequential([model, 
                                         keras.layers.Softmax()])

predictions = probability_model.predict(prediction_pixel_data)

predicted_labels = [np.argmax(probs) for probs in predictions]

predicted_labels_file = open("predicted_labels.dat", "bw")
predicted_labels_file.write(bytearray(predicted_labels))
predicted_labels_file.close()