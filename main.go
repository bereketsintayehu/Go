package main

import (
	"fmt"
)

func main() {
	var (
		studentName string
		numberOfSubjects int
		subjectWithGrade = make(map[string]float64)
		totalGrade float64
	)

	for {
		fmt.Print("Enter your name: ")
		fmt.Scanln(&studentName)
		
		if studentName != "" {
			break
		}
		fmt.Println("Invalid name.")
	}

	for {
		fmt.Print("How many subjects do you have? ")
		fmt.Scanln(&numberOfSubjects)

		if numberOfSubjects > 0 {
			break
		}

		fmt.Println("Invalid number of subjects.")
	}

	for i := 0; i < numberOfSubjects; i++ {
		var subjectName string
		var grade float64

		for {
			fmt.Printf("Enter subject %d name: ", i+1)
			fmt.Scanln(&subjectName)

			if subjectName != "" && subjectWithGrade[subjectName] == 0 {
				break
			}
			fmt.Println("Invalid subject name or subject already entered.")
		}

		for {
			fmt.Printf("Enter %s grade: ", subjectName)
			fmt.Scanln(&grade)

			if grade >= 0 && grade <= 100 {
				break
			}

			fmt.Println("Invalid grade. Grade must be between 0 and 100.")

		}

		subjectWithGrade[subjectName] = grade
		totalGrade += grade
	}
	averageGrade := totalGrade / float64(numberOfSubjects)

	fmt.Printf("Student Name: %s\n", studentName)
	fmt.Println("Your Grades:")
	for subject, grade := range subjectWithGrade {
		fmt.Printf("%s: %.2f\n", subject, grade)
	}
	fmt.Printf("Average Grade: %.2f\n", averageGrade)
}
