package plum

type TrainConfig struct {
	ModelAttributes map[string]string
}

func (t *TrainConfig) TrainAll() {
	for name, model := range App.Models {
		if model.Train != nil {
			t.Train(name)
		}
	}
}

func (t *TrainConfig) Train(modelName string) {
	model := App.Models[modelName]

	App.Log("Model", "Training "+modelName, "orange")
	model.TrainModel(model.SetAttributes(t.ModelAttributes))

}
