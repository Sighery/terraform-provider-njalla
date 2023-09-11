package njalla

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/Sighery/gonjalla"
)

func TestAccServer_Create(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckServerCreate(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckServerExists(
						"njalla_server.test_create",
					),
					resource.TestCheckResourceAttr(
						"njalla_server.test_create",
						"name",
						"test_create",
					),
					resource.TestCheckResourceAttr(
						"njalla_server.test_create",
						"instance_type",
						"njalla1",
					),
					resource.TestCheckResourceAttr(
						"njalla_server.test_create",
						"os",
						"ubuntu2204",
					),
					resource.TestCheckResourceAttr(
						"njalla_server.test_create",
						"public_key",
						"ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQCckbBoD0IkvWOvRHJDTbZgKG70mBh0xvrLZ0BFb6AGzzqCiz9ItqOljtpYp1baAYVygvy147oDR3Z5Zlisl7kJNgPV6YC32ZcFEO/qLczhzeYuH+JzrcCCAsXKdsfVCyPwBmPsI5aL4M+cAgXtii5XoMcKWWr0QUlVUb71gBksr+daLcLSp7F48aIwbM4PM7DaGR6f7b4q/ewmA4QjrUQl6BH/2he6L/gAkVRIOmPXHXoDPgpoSro9xi7NYwMYzetdXqtCj+vSog63NshVvTBNGPFl7TIAteHj4PfuvLqv1t89WFZLVIWuNgjuo4QLS0WlxPs4UiHC2NvGRK7Gju313vYneJH2yNQyX8HZLCAVrtxvCe+rxPD6YqZdkCZLsURG5gK4r8Smc1aDydkPVZCLyvc6C/ZeDFR5snnm8nbm0JHs8g0W/2pEusSOigw+aohDfBopiAFCPnzSZzir3pnbKdd5MMMqN/FfQhQSE1AcjIb3Mkc5akDxJ0dtY89wjMU= yourusername@yourhostname",
					),
					resource.TestCheckResourceAttr(
						"njalla_server.test_create",
						"months",
						"1",
					),
				),
			},
		},
	})
}

func TestAccServer_Update(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckServerCreate(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckServerExists(
						"njalla_server.test_create",
					),
					resource.TestCheckResourceAttr(
						"njalla_server.test_create",
						"name",
						"test_create",
					),
					resource.TestCheckResourceAttr(
						"njalla_server.test_create",
						"instance_type",
						"njalla1",
					),
					resource.TestCheckResourceAttr(
						"njalla_server.test_create",
						"os",
						"ubuntu2204",
					),
					resource.TestCheckResourceAttr(
						"njalla_server.test_create",
						"public_key",
						"ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQCckbBoD0IkvWOvRHJDTbZgKG70mBh0xvrLZ0BFb6AGzzqCiz9ItqOljtpYp1baAYVygvy147oDR3Z5Zlisl7kJNgPV6YC32ZcFEO/qLczhzeYuH+JzrcCCAsXKdsfVCyPwBmPsI5aL4M+cAgXtii5XoMcKWWr0QUlVUb71gBksr+daLcLSp7F48aIwbM4PM7DaGR6f7b4q/ewmA4QjrUQl6BH/2he6L/gAkVRIOmPXHXoDPgpoSro9xi7NYwMYzetdXqtCj+vSog63NshVvTBNGPFl7TIAteHj4PfuvLqv1t89WFZLVIWuNgjuo4QLS0WlxPs4UiHC2NvGRK7Gju313vYneJH2yNQyX8HZLCAVrtxvCe+rxPD6YqZdkCZLsURG5gK4r8Smc1aDydkPVZCLyvc6C/ZeDFR5snnm8nbm0JHs8g0W/2pEusSOigw+aohDfBopiAFCPnzSZzir3pnbKdd5MMMqN/FfQhQSE1AcjIb3Mkc5akDxJ0dtY89wjMU= yourusername@yourhostname",
					),
					resource.TestCheckResourceAttr(
						"njalla_server.test_create",
						"months",
						"1",
					),
				),
			},
			{
				Config: testAccCheckServerUpdatePost(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckServerExists(
						"njalla_server.test_update",
					),
					resource.TestCheckResourceAttr(
						"njalla_server.test_update",
						"name",
						"test_update",
					),
					resource.TestCheckResourceAttr(
						"njalla_server.test_update",
						"instance_type",
						"njalla2",
					),
					resource.TestCheckResourceAttr(
						"njalla_server.test_update",
						"os",
						"ubuntu2004",
					),
					resource.TestCheckResourceAttr(
						"njalla_server.test_update",
						"public_key",
						"ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQCckbBoD0IkvWOvRHJDTbZgKG70mBh0xvrLZ0BFb6AGzzqCiz9ItqOljtpYp1baAYVygvy147oDR3Z5Zlisl7kJNgPV6YC32ZcFEO/qLczhzeYuH+JzrcCCAsXKdsfVCyPwBmPsI5aL4M+cAgXtii5XoMcKWWr0QUlVUb71gBksr+daLcLSp7F48aIwbM4PM7DaGR6f7b4q/ewmA4QjrUQl6BH/2he6L/gAkVRIOmPXHXoDPgpoSro9xi7NYwMYzetdXqtCj+vSog63NshVvTBNGPFl7TIAteHj4PfuvLqv1t89WFZLVIWuNgjuo4QLS0WlxPs4UiHC2NvGRK7Gju313vYneJH2yNQyX8HZLCAVrtxvCe+rxPD6YqZdkCZLsURG5gK4r8Smc1aDydkPVZCLyvc6C/ZeDFR5snnm8nbm0JHs8g0W/2pEusSOigw+aohDfBopiAFCPnzSZzir3pnbKdd5MMMqN/FfQhQSE1AcjIb3Mkc5akDxJ0dtY89wjMU= yourusername@yourhostname2",
					),
					resource.TestCheckResourceAttr(
						"njalla_server.test_update",
						"months",
						"1",
					),
				),
			},
		},
	})
}

func TestAccServer_Import(t *testing.T) {
	
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckServerImport(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckServerExists(
						"njalla_server.test_import",
					),
				),
			},
			{
				ResourceName:        "njalla_server.test_import",
				ImportState:         true,
				ImportStateVerify:   true,
			},
		},
	})
}

func testAccCheckServerCreate() string {
	return fmt.Sprintf(`
resource njalla_server test_create {
	name	=	"test_create"
	instance_type	= "njalla1"
	os	= "ubuntu2204"
	public_key = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQCckbBoD0IkvWOvRHJDTbZgKG70mBh0xvrLZ0BFb6AGzzqCiz9ItqOljtpYp1baAYVygvy147oDR3Z5Zlisl7kJNgPV6YC32ZcFEO/qLczhzeYuH+JzrcCCAsXKdsfVCyPwBmPsI5aL4M+cAgXtii5XoMcKWWr0QUlVUb71gBksr+daLcLSp7F48aIwbM4PM7DaGR6f7b4q/ewmA4QjrUQl6BH/2he6L/gAkVRIOmPXHXoDPgpoSro9xi7NYwMYzetdXqtCj+vSog63NshVvTBNGPFl7TIAteHj4PfuvLqv1t89WFZLVIWuNgjuo4QLS0WlxPs4UiHC2NvGRK7Gju313vYneJH2yNQyX8HZLCAVrtxvCe+rxPD6YqZdkCZLsURG5gK4r8Smc1aDydkPVZCLyvc6C/ZeDFR5snnm8nbm0JHs8g0W/2pEusSOigw+aohDfBopiAFCPnzSZzir3pnbKdd5MMMqN/FfQhQSE1AcjIb3Mkc5akDxJ0dtY89wjMU= yourusername@yourhostname"
	months = 1
}`)
}

func testAccCheckServerUpdatePost() string {
	return fmt.Sprintf(`
resource njalla_server test_update {
	name	=	"test_update"
	instance_type	= "njalla2"
	os	= "ubuntu2004"
	public_key = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQCckbBoD0IkvWOvRHJDTbZgKG70mBh0xvrLZ0BFb6AGzzqCiz9ItqOljtpYp1baAYVygvy147oDR3Z5Zlisl7kJNgPV6YC32ZcFEO/qLczhzeYuH+JzrcCCAsXKdsfVCyPwBmPsI5aL4M+cAgXtii5XoMcKWWr0QUlVUb71gBksr+daLcLSp7F48aIwbM4PM7DaGR6f7b4q/ewmA4QjrUQl6BH/2he6L/gAkVRIOmPXHXoDPgpoSro9xi7NYwMYzetdXqtCj+vSog63NshVvTBNGPFl7TIAteHj4PfuvLqv1t89WFZLVIWuNgjuo4QLS0WlxPs4UiHC2NvGRK7Gju313vYneJH2yNQyX8HZLCAVrtxvCe+rxPD6YqZdkCZLsURG5gK4r8Smc1aDydkPVZCLyvc6C/ZeDFR5snnm8nbm0JHs8g0W/2pEusSOigw+aohDfBopiAFCPnzSZzir3pnbKdd5MMMqN/FfQhQSE1AcjIb3Mkc5akDxJ0dtY89wjMU= yourusername@yourhostname2"
	months = 1
}`)
}

func testAccCheckServerImport() string {
	return fmt.Sprintf(`
resource njalla_server test_import {
	name	=	"test_import"
	instance_type	= "njalla1"
	os	= "ubuntu2204"
	public_key = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQCckbBoD0IkvWOvRHJDTbZgKG70mBh0xvrLZ0BFb6AGzzqCiz9ItqOljtpYp1baAYVygvy147oDR3Z5Zlisl7kJNgPV6YC32ZcFEO/qLczhzeYuH+JzrcCCAsXKdsfVCyPwBmPsI5aL4M+cAgXtii5XoMcKWWr0QUlVUb71gBksr+daLcLSp7F48aIwbM4PM7DaGR6f7b4q/ewmA4QjrUQl6BH/2he6L/gAkVRIOmPXHXoDPgpoSro9xi7NYwMYzetdXqtCj+vSog63NshVvTBNGPFl7TIAteHj4PfuvLqv1t89WFZLVIWuNgjuo4QLS0WlxPs4UiHC2NvGRK7Gju313vYneJH2yNQyX8HZLCAVrtxvCe+rxPD6YqZdkCZLsURG5gK4r8Smc1aDydkPVZCLyvc6C/ZeDFR5snnm8nbm0JHs8g0W/2pEusSOigw+aohDfBopiAFCPnzSZzir3pnbKdd5MMMqN/FfQhQSE1AcjIb3Mkc5akDxJ0dtY89wjMU= yourusername@yourhostname"
	months = 1
}`)
}

func testAccCheckServerExists(resource string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("Not found: %s", resource)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No record ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		servers, err := gonjalla.ListServers(config.Token)
		if err != nil {
			return fmt.Errorf(
				"Error fetching server data: %s",
				err,
			)
		}

		for _, server := range servers {
			if server.ID == rs.Primary.ID {
				return nil
			}
		}

		return fmt.Errorf(
			"The server %s doesn't exist", rs.Primary.ID,
		)
	}
}

func testAccCheckServerDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "njalla_server" {
			continue
		}

		servers, err := gonjalla.ListServers(config.Token)
		if err != nil {
			return fmt.Errorf(
				"Error servers: %s",
				err,
			)
		}

		for _, server := range servers {
			if server.ID == rs.Primary.ID {
				return fmt.Errorf(
					"Server %s still exists",
					rs.Primary.ID,
				)
			}
		}
	}

	return nil
}
