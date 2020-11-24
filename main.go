package main

import (
	"fmt"
	"time"
)

type Proceso struct {
	id uint64
}

func FuncionProceso(id uint64, canalAct chan bool, procesoDetener chan uint64) {
	i := uint64(0)
	for {
		select {
		case detenerProceso := <-procesoDetener:
			if detenerProceso == id {
				return
			}
		case <-canalAct:
			fmt.Printf("id %d: %d\n", id, i)
			i = i + 1
			time.Sleep(time.Millisecond * 500)
		default:
			i = i + 1
			time.Sleep(time.Millisecond * 500)
		}
	}
}

func DetenerProceso(id uint64, procesoDetener chan uint64, procesos []Proceso) {
	for i := 0; i < len(procesos); i++ {
		procesoDetener <- id
	}
	return
}

func Imprimir(canalAct chan bool, noImprimir chan bool) {
	for {
		select {
		case <-noImprimir:
			return
		case <-canalAct:
			canalAct <- true
		default:
			canalAct <- true
		}
	}
}

func BorrarProcesoSlice(id uint64, procesos []Proceso) []Proceso {
	for i, value := range procesos {
		if value.id == id {
			return append(procesos[:i], procesos[i+1:]...)
		}
	}
	return procesos
}

const opcAgregarProceso int = 1
const opcMostrarProcesos int = 2
const opcTerminarProceso int = 3
const opcSalir int = 0

func main() {
	var procesos []Proceso
	var opc int
	var idDetenerProceso uint64
	var input string

	idProceso := uint64(0)
	canalAct := make(chan bool)
	procesoDetener := make(chan uint64)
	noImprimir := make(chan bool)

	fmt.Println("len:",len(procesos))

	for {
		imprimirMenu()
		fmt.Scanln(&opc)

		if opc == opcSalir {
			fmt.Println("\nSalir del administrador de procesos")
			fmt.Printf("\n")
			break
		}

		switch opc {
		case opcAgregarProceso:
			procesoNuevo := new(Proceso)
			procesoNuevo.id = idProceso
			fmt.Println("Se ha agregado el proceso con el id: ", idProceso)
			idProceso++
			procesos = append(procesos, *procesoNuevo)

			go FuncionProceso(procesoNuevo.id, canalAct, procesoDetener)

		case opcMostrarProcesos:
			if len(procesos) > 0{
				canalAct <- true
				go Imprimir(canalAct, noImprimir)

				fmt.Scanln(&input)
				noImprimir <- true
			} else {
				fmt.Printf("\nNo se ha agregado un proceso")
			}

		case opcTerminarProceso:
			fmt.Printf("\nIngrese el ID del proceso que desea detener: ")
			fmt.Scanln(&idDetenerProceso)
			
			go DetenerProceso(idDetenerProceso, procesoDetener, procesos)
			fmt.Println("\nSe ha eliminado el proceso con el id: ", idDetenerProceso)

			procesos = BorrarProcesoSlice(idDetenerProceso, procesos)

			fmt.Scanln(&input)
		}
	}
}

func imprimirMenu(){
	fmt.Printf("\n\t***Administrador de Procesos***\n\n")
	fmt.Println("1) Agregar proceso")
	fmt.Println("2) Mostrar procesos")
	fmt.Println("3) Terminar proceso")
	fmt.Println("0) Salir")
	fmt.Printf("Ingrese una opción: ")
}