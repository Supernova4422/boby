package command

type OxfordTranslateResponseStruct struct {
	Metadata struct {
	} `json:"metadata"`
	Results []struct {
		ID             string `json:"id"`
		Language       string `json:"language"`
		LexicalEntries []struct {
			Compounds []struct {
				Domains []struct {
					ID   string `json:"id"`
					Text string `json:"text"`
				} `json:"domains"`
				ID       string `json:"id"`
				Language string `json:"language"`
				Regions  []struct {
					ID   string `json:"id"`
					Text string `json:"text"`
				} `json:"regions"`
				Registers []struct {
					ID   string `json:"id"`
					Text string `json:"text"`
				} `json:"registers"`
				Text string `json:"text"`
			} `json:"compounds"`
			DerivativeOf []struct {
				Domains []struct {
					ID   string `json:"id"`
					Text string `json:"text"`
				} `json:"domains"`
				ID       string `json:"id"`
				Language string `json:"language"`
				Regions  []struct {
					ID   string `json:"id"`
					Text string `json:"text"`
				} `json:"regions"`
				Registers []struct {
					ID   string `json:"id"`
					Text string `json:"text"`
				} `json:"registers"`
				Text string `json:"text"`
			} `json:"derivativeOf"`
			Derivatives []struct {
				Domains []struct {
					ID   string `json:"id"`
					Text string `json:"text"`
				} `json:"domains"`
				ID       string `json:"id"`
				Language string `json:"language"`
				Regions  []struct {
					ID   string `json:"id"`
					Text string `json:"text"`
				} `json:"regions"`
				Registers []struct {
					ID   string `json:"id"`
					Text string `json:"text"`
				} `json:"registers"`
				Text string `json:"text"`
			} `json:"derivatives"`
			Entries []struct {
				CrossReferenceMarkers []string `json:"crossReferenceMarkers"`
				CrossReferences       []struct {
					ID   string `json:"id"`
					Text string `json:"text"`
					Type string `json:"type"`
				} `json:"crossReferences"`
				Etymologies         []string `json:"etymologies"`
				GrammaticalFeatures []struct {
					ID   string `json:"id"`
					Text string `json:"text"`
					Type string `json:"type"`
				} `json:"grammaticalFeatures"`
				HomographNumber string `json:"homographNumber"`
				Inflections     []struct {
					Domains []struct {
						ID   string `json:"id"`
						Text string `json:"text"`
					} `json:"domains"`
					GrammaticalFeatures []struct {
						ID   string `json:"id"`
						Text string `json:"text"`
						Type string `json:"type"`
					} `json:"grammaticalFeatures"`
					InflectedForm   string `json:"inflectedForm"`
					LexicalCategory struct {
						ID   string `json:"id"`
						Text string `json:"text"`
					} `json:"lexicalCategory"`
					Pronunciations []struct {
						AudioFile        string   `json:"audioFile"`
						Dialects         []string `json:"dialects"`
						PhoneticNotation string   `json:"phoneticNotation"`
						PhoneticSpelling string   `json:"phoneticSpelling"`
						Regions          []struct {
							ID   string `json:"id"`
							Text string `json:"text"`
						} `json:"regions"`
						Registers []struct {
							ID   string `json:"id"`
							Text string `json:"text"`
						} `json:"registers"`
					} `json:"pronunciations"`
					Regions []struct {
						ID   string `json:"id"`
						Text string `json:"text"`
					} `json:"regions"`
					Registers []struct {
						ID   string `json:"id"`
						Text string `json:"text"`
					} `json:"registers"`
				} `json:"inflections"`
				Notes []struct {
					ID   string `json:"id"`
					Text string `json:"text"`
					Type string `json:"type"`
				} `json:"notes"`
				Pronunciations []struct {
					AudioFile        string   `json:"audioFile"`
					Dialects         []string `json:"dialects"`
					PhoneticNotation string   `json:"phoneticNotation"`
					PhoneticSpelling string   `json:"phoneticSpelling"`
					Regions          []struct {
						ID   string `json:"id"`
						Text string `json:"text"`
					} `json:"regions"`
					Registers []struct {
						ID   string `json:"id"`
						Text string `json:"text"`
					} `json:"registers"`
				} `json:"pronunciations"`
				Senses []struct {
					Antonyms []struct {
						Domains []struct {
							ID   string `json:"id"`
							Text string `json:"text"`
						} `json:"domains"`
						ID       string `json:"id"`
						Language string `json:"language"`
						Regions  []struct {
							ID   string `json:"id"`
							Text string `json:"text"`
						} `json:"regions"`
						Registers []struct {
							ID   string `json:"id"`
							Text string `json:"text"`
						} `json:"registers"`
						Text string `json:"text"`
					} `json:"antonyms"`
					Constructions []struct {
						Domains []struct {
							ID   string `json:"id"`
							Text string `json:"text"`
						} `json:"domains"`
						Examples [][]string `json:"examples"`
						Notes    []struct {
							ID   string `json:"id"`
							Text string `json:"text"`
							Type string `json:"type"`
						} `json:"notes"`
						Regions []struct {
							ID   string `json:"id"`
							Text string `json:"text"`
						} `json:"regions"`
						Registers []struct {
							ID   string `json:"id"`
							Text string `json:"text"`
						} `json:"registers"`
						Text         string `json:"text"`
						Translations []struct {
							Collocations []struct {
								ID   string `json:"id"`
								Text string `json:"text"`
								Type string `json:"type"`
							} `json:"collocations"`
							Domains []struct {
								ID   string `json:"id"`
								Text string `json:"text"`
							} `json:"domains"`
							GrammaticalFeatures []struct {
								ID   string `json:"id"`
								Text string `json:"text"`
								Type string `json:"type"`
							} `json:"grammaticalFeatures"`
							Language string `json:"language"`
							Notes    []struct {
								ID   string `json:"id"`
								Text string `json:"text"`
								Type string `json:"type"`
							} `json:"notes"`
							Regions []struct {
								ID   string `json:"id"`
								Text string `json:"text"`
							} `json:"regions"`
							Registers []struct {
								ID   string `json:"id"`
								Text string `json:"text"`
							} `json:"registers"`
							Text       string `json:"text"`
							ToneGroups []struct {
								Tones []struct {
									Type  string `json:"type"`
									Value string `json:"value"`
								} `json:"tones"`
							} `json:"toneGroups"`
							Type string `json:"type"`
						} `json:"translations"`
					} `json:"constructions"`
					CrossReferenceMarkers []string `json:"crossReferenceMarkers"`
					CrossReferences       []struct {
						ID   string `json:"id"`
						Text string `json:"text"`
						Type string `json:"type"`
					} `json:"crossReferences"`
					DatasetCrossLinks []struct {
						Language string `json:"language"`
						EntryID  string `json:"entry_id"`
						SenseID  string `json:"sense_id"`
					} `json:"datasetCrossLinks"`
					Definitions   []string `json:"definitions"`
					DomainClasses []struct {
						ID   string `json:"id"`
						Text string `json:"text"`
					} `json:"domainClasses"`
					Domains []struct {
						ID   string `json:"id"`
						Text string `json:"text"`
					} `json:"domains"`
					Etymologies []string `json:"etymologies"`
					Examples    []struct {
						Collocations []struct {
							ID   string `json:"id"`
							Text string `json:"text"`
							Type string `json:"type"`
						} `json:"collocations"`
						CrossReferenceMarkers []string `json:"crossReferenceMarkers"`
						CrossReferences       []struct {
							ID   string `json:"id"`
							Text string `json:"text"`
							Type string `json:"type"`
						} `json:"crossReferences"`
						Definitions []string `json:"definitions"`
						Domains     []struct {
							ID   string `json:"id"`
							Text string `json:"text"`
						} `json:"domains"`
						Notes []struct {
							ID   string `json:"id"`
							Text string `json:"text"`
							Type string `json:"type"`
						} `json:"notes"`
						Regions []struct {
							ID   string `json:"id"`
							Text string `json:"text"`
						} `json:"regions"`
						Registers []struct {
							ID   string `json:"id"`
							Text string `json:"text"`
						} `json:"registers"`
						SenseIds     []string `json:"senseIds"`
						Text         string   `json:"text"`
						Translations []struct {
							Collocations []struct {
								ID   string `json:"id"`
								Text string `json:"text"`
								Type string `json:"type"`
							} `json:"collocations"`
							Domains []struct {
								ID   string `json:"id"`
								Text string `json:"text"`
							} `json:"domains"`
							GrammaticalFeatures []struct {
								ID   string `json:"id"`
								Text string `json:"text"`
								Type string `json:"type"`
							} `json:"grammaticalFeatures"`
							Language string `json:"language"`
							Notes    []struct {
								ID   string `json:"id"`
								Text string `json:"text"`
								Type string `json:"type"`
							} `json:"notes"`
							Regions []struct {
								ID   string `json:"id"`
								Text string `json:"text"`
							} `json:"regions"`
							Registers []struct {
								ID   string `json:"id"`
								Text string `json:"text"`
							} `json:"registers"`
							Text       string `json:"text"`
							ToneGroups []struct {
								Tones []struct {
									Type  string `json:"type"`
									Value string `json:"value"`
								} `json:"tones"`
							} `json:"toneGroups"`
							Type string `json:"type"`
						} `json:"translations"`
					} `json:"examples"`
					ID          string `json:"id"`
					Inflections []struct {
						Domains []struct {
							ID   string `json:"id"`
							Text string `json:"text"`
						} `json:"domains"`
						GrammaticalFeatures []struct {
							ID   string `json:"id"`
							Text string `json:"text"`
							Type string `json:"type"`
						} `json:"grammaticalFeatures"`
						InflectedForm   string `json:"inflectedForm"`
						LexicalCategory struct {
							ID   string `json:"id"`
							Text string `json:"text"`
						} `json:"lexicalCategory"`
						Pronunciations []struct {
							AudioFile        string   `json:"audioFile"`
							Dialects         []string `json:"dialects"`
							PhoneticNotation string   `json:"phoneticNotation"`
							PhoneticSpelling string   `json:"phoneticSpelling"`
							Regions          []struct {
								ID   string `json:"id"`
								Text string `json:"text"`
							} `json:"regions"`
							Registers []struct {
								ID   string `json:"id"`
								Text string `json:"text"`
							} `json:"registers"`
						} `json:"pronunciations"`
						Regions []struct {
							ID   string `json:"id"`
							Text string `json:"text"`
						} `json:"regions"`
						Registers []struct {
							ID   string `json:"id"`
							Text string `json:"text"`
						} `json:"registers"`
					} `json:"inflections"`
					Notes []struct {
						ID   string `json:"id"`
						Text string `json:"text"`
						Type string `json:"type"`
					} `json:"notes"`
					Pronunciations []struct {
						AudioFile        string   `json:"audioFile"`
						Dialects         []string `json:"dialects"`
						PhoneticNotation string   `json:"phoneticNotation"`
						PhoneticSpelling string   `json:"phoneticSpelling"`
						Regions          []struct {
							ID   string `json:"id"`
							Text string `json:"text"`
						} `json:"regions"`
						Registers []struct {
							ID   string `json:"id"`
							Text string `json:"text"`
						} `json:"registers"`
					} `json:"pronunciations"`
					Regions []struct {
						ID   string `json:"id"`
						Text string `json:"text"`
					} `json:"regions"`
					Registers []struct {
						ID   string `json:"id"`
						Text string `json:"text"`
					} `json:"registers"`
					SemanticClasses []struct {
						ID   string `json:"id"`
						Text string `json:"text"`
					} `json:"semanticClasses"`
					ShortDefinitions []string `json:"shortDefinitions"`
					Subsenses        []struct {
					} `json:"subsenses"`
					Synonyms []struct {
						Domains []struct {
							ID   string `json:"id"`
							Text string `json:"text"`
						} `json:"domains"`
						ID       string `json:"id"`
						Language string `json:"language"`
						Regions  []struct {
							ID   string `json:"id"`
							Text string `json:"text"`
						} `json:"regions"`
						Registers []struct {
							ID   string `json:"id"`
							Text string `json:"text"`
						} `json:"registers"`
						Text string `json:"text"`
					} `json:"synonyms"`
					ThesaurusLinks []struct {
						EntryID string `json:"entry_id"`
						SenseID string `json:"sense_id"`
					} `json:"thesaurusLinks"`
					Translations []struct {
						Collocations []struct {
							ID   string `json:"id"`
							Text string `json:"text"`
							Type string `json:"type"`
						} `json:"collocations"`
						Domains []struct {
							ID   string `json:"id"`
							Text string `json:"text"`
						} `json:"domains"`
						GrammaticalFeatures []struct {
							ID   string `json:"id"`
							Text string `json:"text"`
							Type string `json:"type"`
						} `json:"grammaticalFeatures"`
						Language string `json:"language"`
						Notes    []struct {
							ID   string `json:"id"`
							Text string `json:"text"`
							Type string `json:"type"`
						} `json:"notes"`
						Regions []struct {
							ID   string `json:"id"`
							Text string `json:"text"`
						} `json:"regions"`
						Registers []struct {
							ID   string `json:"id"`
							Text string `json:"text"`
						} `json:"registers"`
						Text       string `json:"text"`
						ToneGroups []struct {
							Tones []struct {
								Type  string `json:"type"`
								Value string `json:"value"`
							} `json:"tones"`
						} `json:"toneGroups"`
						Type string `json:"type"`
					} `json:"translations"`
					VariantForms []struct {
						Domains []struct {
							ID   string `json:"id"`
							Text string `json:"text"`
						} `json:"domains"`
						Notes []struct {
							ID   string `json:"id"`
							Text string `json:"text"`
							Type string `json:"type"`
						} `json:"notes"`
						Pronunciations []struct {
							AudioFile        string   `json:"audioFile"`
							Dialects         []string `json:"dialects"`
							PhoneticNotation string   `json:"phoneticNotation"`
							PhoneticSpelling string   `json:"phoneticSpelling"`
							Regions          []struct {
								ID   string `json:"id"`
								Text string `json:"text"`
							} `json:"regions"`
							Registers []struct {
								ID   string `json:"id"`
								Text string `json:"text"`
							} `json:"registers"`
						} `json:"pronunciations"`
						Regions []struct {
							ID   string `json:"id"`
							Text string `json:"text"`
						} `json:"regions"`
						Registers []struct {
							ID   string `json:"id"`
							Text string `json:"text"`
						} `json:"registers"`
						Text string `json:"text"`
					} `json:"variantForms"`
				} `json:"senses"`
				VariantForms []struct {
					Domains []struct {
						ID   string `json:"id"`
						Text string `json:"text"`
					} `json:"domains"`
					Notes []struct {
						ID   string `json:"id"`
						Text string `json:"text"`
						Type string `json:"type"`
					} `json:"notes"`
					Pronunciations []struct {
						AudioFile        string   `json:"audioFile"`
						Dialects         []string `json:"dialects"`
						PhoneticNotation string   `json:"phoneticNotation"`
						PhoneticSpelling string   `json:"phoneticSpelling"`
						Regions          []struct {
							ID   string `json:"id"`
							Text string `json:"text"`
						} `json:"regions"`
						Registers []struct {
							ID   string `json:"id"`
							Text string `json:"text"`
						} `json:"registers"`
					} `json:"pronunciations"`
					Regions []struct {
						ID   string `json:"id"`
						Text string `json:"text"`
					} `json:"regions"`
					Registers []struct {
						ID   string `json:"id"`
						Text string `json:"text"`
					} `json:"registers"`
					Text string `json:"text"`
				} `json:"variantForms"`
			} `json:"entries"`
			GrammaticalFeatures []struct {
				ID   string `json:"id"`
				Text string `json:"text"`
				Type string `json:"type"`
			} `json:"grammaticalFeatures"`
			Language        string `json:"language"`
			LexicalCategory struct {
				ID   string `json:"id"`
				Text string `json:"text"`
			} `json:"lexicalCategory"`
			Notes []struct {
				ID   string `json:"id"`
				Text string `json:"text"`
				Type string `json:"type"`
			} `json:"notes"`
			PhrasalVerbs []struct {
				Domains []struct {
					ID   string `json:"id"`
					Text string `json:"text"`
				} `json:"domains"`
				ID       string `json:"id"`
				Language string `json:"language"`
				Regions  []struct {
					ID   string `json:"id"`
					Text string `json:"text"`
				} `json:"regions"`
				Registers []struct {
					ID   string `json:"id"`
					Text string `json:"text"`
				} `json:"registers"`
				Text string `json:"text"`
			} `json:"phrasalVerbs"`
			Phrases []struct {
				Domains []struct {
					ID   string `json:"id"`
					Text string `json:"text"`
				} `json:"domains"`
				ID       string `json:"id"`
				Language string `json:"language"`
				Regions  []struct {
					ID   string `json:"id"`
					Text string `json:"text"`
				} `json:"regions"`
				Registers []struct {
					ID   string `json:"id"`
					Text string `json:"text"`
				} `json:"registers"`
				Text string `json:"text"`
			} `json:"phrases"`
			Pronunciations []struct {
				AudioFile        string   `json:"audioFile"`
				Dialects         []string `json:"dialects"`
				PhoneticNotation string   `json:"phoneticNotation"`
				PhoneticSpelling string   `json:"phoneticSpelling"`
				Regions          []struct {
					ID   string `json:"id"`
					Text string `json:"text"`
				} `json:"regions"`
				Registers []struct {
					ID   string `json:"id"`
					Text string `json:"text"`
				} `json:"registers"`
			} `json:"pronunciations"`
			Root         string `json:"root"`
			Text         string `json:"text"`
			VariantForms []struct {
				Domains []struct {
					ID   string `json:"id"`
					Text string `json:"text"`
				} `json:"domains"`
				Notes []struct {
					ID   string `json:"id"`
					Text string `json:"text"`
					Type string `json:"type"`
				} `json:"notes"`
				Pronunciations []struct {
					AudioFile        string   `json:"audioFile"`
					Dialects         []string `json:"dialects"`
					PhoneticNotation string   `json:"phoneticNotation"`
					PhoneticSpelling string   `json:"phoneticSpelling"`
					Regions          []struct {
						ID   string `json:"id"`
						Text string `json:"text"`
					} `json:"regions"`
					Registers []struct {
						ID   string `json:"id"`
						Text string `json:"text"`
					} `json:"registers"`
				} `json:"pronunciations"`
				Regions []struct {
					ID   string `json:"id"`
					Text string `json:"text"`
				} `json:"regions"`
				Registers []struct {
					ID   string `json:"id"`
					Text string `json:"text"`
				} `json:"registers"`
				Text string `json:"text"`
			} `json:"variantForms"`
		} `json:"lexicalEntries"`
		Pronunciations []struct {
			AudioFile        string   `json:"audioFile"`
			Dialects         []string `json:"dialects"`
			PhoneticNotation string   `json:"phoneticNotation"`
			PhoneticSpelling string   `json:"phoneticSpelling"`
			Regions          []struct {
				ID   string `json:"id"`
				Text string `json:"text"`
			} `json:"regions"`
			Registers []struct {
				ID   string `json:"id"`
				Text string `json:"text"`
			} `json:"registers"`
		} `json:"pronunciations"`
		Type string `json:"type"`
		Word string `json:"word"`
	} `json:"results"`
}
