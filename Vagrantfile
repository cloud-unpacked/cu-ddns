# -*- mode: ruby -*-
# vi: set ft=ruby :

Vagrant.configure("2") do |config|
	config.vm.box = "bento/ubuntu-18.04"

	config.vm.provision :shell, path: "vagrant-bootstrap.sh"

	config.vm.provider "virtualbox" do |vb|
		vb.memory = "2048"
		vb.customize [ "modifyvm", :id, "--uartmode1", "disconnected" ]
	end
end
