import matplotlib.pyplot as plt
import numpy as np
import keras
import gzip

# Kind of like how much "context" each pixel gets

training_matrix_size: int
prediction_matrix_size: int

def load_training_pixel_data():
    compressed_pixels_file = open("training_pixel_data.tpixeldata", "br")
    pixels_file = gzip.decompress(compressed_pixels_file.read())
    global training_matrix_size
    training_matrix_size = int.from_bytes(pixels_file[0:2], "big")
    pixels = pixels_file[2:]
    tensors = [np.array(bytearray(pixels[i:i+training_matrix_size*3])) for i in range(0, len(pixels), training_matrix_size*3)]
    return np.array(tensors)

def load_prediction_pixel_data():
    compressed_pixels_file = open("prediction_pixel_data.tpixeldata", "br")
    pixels_file = gzip.decompress(compressed_pixels_file.read())
    global prediction_matrix_size
    prediction_matrix_size = int.from_bytes(pixels_file[0:2], "big")
    pixels = pixels_file[2:]
    tensors = [np.array(bytearray(pixels[i:i+prediction_matrix_size*3])) for i in range(0, len(pixels), prediction_matrix_size*3)]
    return np.array(tensors)

def load_training_label_data():
    labels_file = open("training_label_data.tlabeldata", "br")
    labels_bytes = labels_file.read()
    return np.array(list(labels_bytes))

    
print("Processing pixel training data...")
training_pixel_data = load_training_pixel_data()
print(training_pixel_data.shape)
print("Processing training label data...")
training_label_data = load_training_label_data()
print("Processing prediction data...")
prediction_pixel_data = load_prediction_pixel_data()

if training_matrix_size != prediction_matrix_size:
    print("Matrix sizes are not equal")
    exit(1)

training_pixel_data = training_pixel_data / 255.0

model = keras.Sequential([
    keras.layers.Dense(256, activation='relu'),
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

model.save("boardnet.keras")